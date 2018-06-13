CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/aaronland/go-feed-reader
	cp -r assets src/github.com/aaronland/go-feed-reader/
	cp -r crumb src/github.com/aaronland/go-feed-reader/
	cp -r login src/github.com/aaronland/go-feed-reader/
	cp -r http src/github.com/aaronland/go-feed-reader/
	cp -r password src/github.com/aaronland/go-feed-reader/
	cp -r tables src/github.com/aaronland/go-feed-reader/
	cp -r user src/github.com/aaronland/go-feed-reader/
	cp *.go src/github.com/aaronland/go-feed-reader/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/zendesk/go-bindata/"
	@GOPATH=$(GOPATH) go get -u "github.com/aaronland/go-secretbox"
	@GOPATH=$(GOPATH) go get -u "github.com/aaronland/go-string"
	@GOPATH=$(GOPATH) go get -u "github.com/aaronland/go-sql-pagination"
	@GOPATH=$(GOPATH) go get -u "github.com/arschles/go-bindata-html-template"
	@GOPATH=$(GOPATH) go get -u "github.com/mmcdole/gofeed"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-sanitize"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-sqlite"
	@GOPATH=$(GOPATH) go get -u "github.com/grokify/html-strip-tags-go"
	@GOPATH=$(GOPATH) go get -u "github.com/patrickmn/go-hmaccrypt"
	@GOPATH=$(GOPATH) go get -u "github.com/Pallinder/go-randomdata"

	rm -rf src/github.com/zendesk/go-bindata/testdata

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt crumb/*.go
	go fmt http/*.go
	go fmt login/*.go
	go fmt password/*.go
	go fmt tables/*.go
	go fmt user/*.go
	go fmt *.go

assets: self
	@GOPATH=$(GOPATH) go build -o bin/go-bindata ./vendor/github.com/zendesk/go-bindata/go-bindata/
	rm -rf templates/*/*~
	rm -rf assets
	mkdir -p assets/html
	@GOPATH=$(GOPATH) bin/go-bindata -pkg html -o assets/html/html.go templates/html templates/atom

bin: 	self
	rm -rf bin/*
	# @GOPATH=$(GOPATH) go build --tags "json1 fts5" -o bin/fr-dump cmd/fr-dump.go
	# @GOPATH=$(GOPATH) go build --tags "json1 fts5" -o bin/fr-add cmd/fr-add.go
	# @GOPATH=$(GOPATH) go build --tags "json1 fts5" -o bin/fr-search cmd/fr-search.go
	@GOPATH=$(GOPATH) go build --tags "json1 fts5" -o bin/fr-refresh cmd/fr-refresh.go
	@GOPATH=$(GOPATH) go build --tags "json1 fts5" -o bin/fr-server cmd/fr-server.go

