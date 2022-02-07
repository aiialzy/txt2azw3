package main

import (
	"os"
	"path"
	"text/template"
)

func genInnerTitlePage(novel Novel) {
	t, err := template.ParseFS(root, "templates/inner_title_page.tmpl")
	if err != nil {
		panic(err)
	}

	args := map[string]interface{}{
		"title":  novel.Title,
		"author": novel.Author,
	}

	f, err := os.OpenFile(path.Join(novel.Title, "text", "internal_titlepage.xhtml"), os.O_CREATE|os.O_WRONLY, 0766)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	err = t.Execute(f, args)
	if err != nil {
		panic(err)
	}
}
