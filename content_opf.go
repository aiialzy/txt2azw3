package main

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

type OpfItem struct {
	Id   string
	Href string
}

func genContentOpf(novel Novel) {

	t, err := template.ParseFiles("templates/content.opf.tmpl")
	if err != nil {
		panic(err)
	}

	filePath := path.Join(novel.Title, "content.opf")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0766)
	if err != nil {
		panic(err)
	}

	var items []OpfItem
	for _, chapter := range novel.Chapters {
		item := OpfItem{
			Id:   fmt.Sprintf("chapter%d", chapter.Index),
			Href: path.Join("text", chapter.Filename),
		}
		items = append(items, item)
	}

	args := map[string]interface{}{
		"title":    novel.Title,
		"author":   novel.Author,
		"uuid":     novel.UUID,
		"items":    items,
		"language": novel.Language,
	}

	err = t.Execute(f, args)
	if err != nil {
		panic(err)
	}
}
