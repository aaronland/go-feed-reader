package http

import (
	"errors"
	"github.com/aaronland/go-feed-reader/assets/html"
	"html/template"
	"io"
	"log"
	"path/filepath"
)

type Templates struct {
	lookup map[string]*template.Template
}

func CompileTemplates(paths []string) (*Templates, error) {

	lookup := make(map[string]*template.Template)

	for _, path := range paths {

		fname := filepath.Base(path)

		body, err := html.Asset(path)

		if err != nil {
			return nil, err
		}

		log.Println("ADD", fname)

		t, err := template.New(fname).Parse(string(body))

		if err != nil {
			return nil, err
		}

		lookup[fname] = t
	}

	t := Templates{lookup}
	return &t, nil
}

func (t *Templates) ExecuteTemplate(wr io.Writer, name string, vars interface{}) error {

	tpl, ok := t.lookup[name]

	if !ok {
		return errors.New("Invalid or unknown template name")
	}

	return tpl.Execute(wr, vars)
}
