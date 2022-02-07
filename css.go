package main

import (
	_ "embed"
	"io/ioutil"
	"path"
)

//go:embed templates/css.tmpl
var css string

func genCss(novel Novel) {

	err := ioutil.WriteFile(path.Join(novel.Title, "styles", "stylesheet.css"), []byte(css), 0766)
	if err != nil {
		panic(err)
	}
}
