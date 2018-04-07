CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/aaronland/go-feed-reader
	cp -r assets src/github.com/aaronland/go-feed-reader/
	cp -r http src/github.com/aaronland/go-feed-reader/
	cp -r tables src/github.com/aaronland/go-feed-reader/
	cp *.go src/github.com/aaronland/go-feed-reader/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/zendesk/go-bindata/"
	@GOPATH=$(GOPATH) go get -u "github.com/mmcdole/gofeed"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-sqlite"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt tables/*.go
	go fmt http/*.go
	go fmt *.go

assets: self
	@GOPATH=$(GOPATH) go build -o bin/go-bindata ./vendor/github.com/zendesk/go-bindata/go-bindata/
	rm -rf templates/*/*~
	rm -rf assets
	mkdir -p assets/html
	@GOPATH=$(GOPATH) bin/go-bindata -pkg html -o assets/html/html.go templates/html

bin: 	self
	@GOPATH=$(GOPATH) go build --tags "json1" -o bin/feed cmd/feed.go
	@GOPATH=$(GOPATH) go build --tags "json1" -o bin/server cmd/server.go

