package main

import (
	"os"
	"path"
	"text/template"

	"github.com/google/uuid"
)

type NavMapItem struct {
	Id        string
	PlayOrder int
	Title     string
	Filename  string
}

func genTocNcx(novel Novel) {
	t, err := template.ParseFS(root, "templates/toc.ncx.tmpl")
	if err != nil {
		panic(err)
	}

	var navmap []NavMapItem

	for _, chapter := range novel.Chapters {
		id, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		item := NavMapItem{
			Id:        id.String(),
			PlayOrder: chapter.Index + 1,
			Title:     chapter.Title,
			Filename:  chapter.Filename,
		}
		navmap = append(navmap, item)
	}

	args := map[string]interface{}{
		"uuid":   novel.UUID,
		"title":  novel.Title,
		"navmap": navmap,
	}

	f, err := os.OpenFile(path.Join(novel.Title, "toc.ncx"), os.O_CREATE|os.O_WRONLY, 0766)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = t.Execute(f, args)
	if err != nil {
		panic(err)
	}
}
