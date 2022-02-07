package main

import (
	"io"
	"os"
	"path"
)

func genGuide(novel Novel) {

	dest, err := os.OpenFile(path.Join(novel.Title, "text", "titlepage.xhtml"), os.O_CREATE|os.O_WRONLY, 0766)

	if err != nil {
		panic(err)
	}

	src, err := root.Open("templates/guide.tmpl")

	if err != nil {
		panic(err)
	}

	_, err = io.Copy(dest, src)

	if err != nil {
		panic(err)
	}

}
