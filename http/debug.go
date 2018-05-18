package http

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/aaronland/go-feed-reader"
	_ "log"
	gohttp "net/http"
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

func entry() Entry {

	dt := randomdata.FullDate()

	e := Entry{
		Id:        "",
		Title:     randomdata.SillyName(),
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

			e := entry()
			entries = append(entries, e)
		}

		author := Author{
			Name: randomdata.FullName(randomdata.RandomGender),
			URI:  "",
		}

		f := Feed{
			Title:    "",
			Subtitle: "",
			Id:       "",
			Updated:  "",
			Link:     "",
			Author:   author,
			Entries:  entries,
		}

		vars := DebugVars{
			Feed: f,
		}

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
