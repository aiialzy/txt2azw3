package main

import (
	"io/ioutil"
	"path"
)

const css = `body {
  display: block;
  margin: 5pt;
  page-break-before: always;
  text-align: justify;
}
h1, h2, h3 {
  font-weight: bold;
  margin-bottom: 1em;
  margin-left: 0;
  margin-right: 0;
  margin-top: 1em;
}
p {
  margin-bottom: 1em;
  margin-left: 0;
  margin-right: 0;
  margin-top: 1em;
}
a {
  color: inherit;
  text-decoration: inherit;
  cursor: default;
}
a[href] {
  color: blue;
  text-decoration: none;
  cursor: pointer;
}
a[href]:hover {
  color: red;
}
code, pre {
  white-space: pre-wrap;
}
.center {
  text-align: center;
}
.cover {
  height: 100%;
}`

func genCss(novel Novel) {

	err := ioutil.WriteFile(path.Join(novel.Title, "styles", "stylesheet.css"), []byte(css), 0766)
	if err != nil {
		panic(err)
	}
}
