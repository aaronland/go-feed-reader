package http

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/aaronland/go-feed-reader"
	_ "log"
	gohttp "net/http"
	"strings"
	"time"
)

type Feed struct {
	Title    string
	Subtitle string
	Id       string
	Updated  string
	Link     string
	Author   Author
	Entries  []Entry
}

type Author struct {
	Name string
	URI  string
}

type Entry struct {
	Id        string
	Link      string
	Title     string
	Published string
	Updated   string
	Content   string
}

type DebugVars struct {
	Feed Feed
}

func entry(req *gohttp.Request) Entry {

	now := time.Now()
	dt := now.Format(time.RFC3339)
	ts := now.Unix()

	title := randomdata.SillyName()
	id := strings.ToLower(title)

	e := Entry{
		Id:        fmt.Sprintf("x-urn-debug-%d#%s", ts, id),
		Link:      fmt.Sprintf("//%s/debug#%s", req.Host, id),
		Title:     title,
		Published: dt,
		Updated:   dt,
		Content:   randomdata.Paragraph(),
	}

	return e
}

func DebugHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	files := []string{
		"templates/atom/debug.xml",
	}

	t, err := CompileTemplate("feed", files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		entries := make([]Entry, 0)

		for len(entries) < 15 {

			e := entry(req)
			entries = append(entries, e)
		}

		author := Author{
			Name: randomdata.FullName(randomdata.RandomGender),
			URI:  fmt.Sprintf("//%s", req.Host),
		}

		ts := time.Now()

		f := Feed{
			Title:    randomdata.SillyName(),
			Subtitle: randomdata.SillyName(),
			Id:       fmt.Sprintf("//%s/debug", req.Host),
			Updated:  ts.Format(time.RFC3339),
			Link:     fmt.Sprintf("//%s/debug", req.Host),
			Author:   author,
			Entries:  entries,
		}

		vars := DebugVars{
			Feed: f,
		}

		rsp.Header().Set("Content-Type", "application/atom+xml;charset=utf-8")

		err = t.Execute(rsp, vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
