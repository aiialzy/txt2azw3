package main

import (
	"io/ioutil"
	"path"
)

func genGuide(novel Novel) {
	src := `<?xml version='1.0' encoding='utf-8'?>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">

<head>
  <title>Cover</title>
  <style type="text/css">
  @page {
    padding: 0;
    margin: 0;
    }
  body {
    text-align: center;
    padding: 0;
    margin: 0;
  }
  </style>
</head>

<body>

  <div>

    <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" height="100%" preserveAspectRatio="none" version="1.1" viewBox="0 0 600 800" width="100%">
      <image height="800" width="600" xlink:href="../images/cover.jpg"/>
    </svg>

  </div>

</body>

</html>`

	err := ioutil.WriteFile(path.Join(novel.Title, "text", "titlepage.xhtml"), []byte(src), 0766)
	if err != nil {
		panic(err)
	}

}
