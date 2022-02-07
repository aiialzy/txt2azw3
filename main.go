package main

import (
	"embed"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
)

//go:embed templates/**
//go:embed font.ttf
var root embed.FS

type Novel struct {
	TxtPath  string
	UUID     string
	Title    string
	Author   string
	Language string
	Chapters []Chapter
}

func genBaseDir(novel Novel) {
	dirs := []string{"images", "META-INF", "styles", "text"}
	for _, dir := range dirs {
		err := os.MkdirAll(path.Join(novel.Title, dir), 0766)
		if err != nil {
			panic(err)
		}
	}

}

func genMimetype(novel Novel) {
	err := ioutil.WriteFile(path.Join(novel.Title, "mimetype"), []byte("application/epub+zip"), 0766)
	if err != nil {
		panic(err)
	}
}

func genContainerXML(novel Novel) {
	containerXML := `<?xml version='1.0' encoding='utf-8'?>
<container xmlns="urn:oasis:names:tc:opendocument:xmlns:container" version="1.0">
    <rootfiles>
        <rootfile full-path="content.opf" media-type="application/oebps-package+xml"/>
   </rootfiles>
</container>`

	err := ioutil.WriteFile(path.Join(novel.Title, "META-INF", "container.xml"), []byte(containerXML), 0766)
	if err != nil {
		panic(err)
	}
}

func clean(novel Novel) {
	files := []string{novel.Title, "out.txt"}
	for _, f := range files {
		if err := os.RemoveAll(f); err != nil {
			log.Println(err)
		}
	}
}

func createNovel(txtPath, title, author, language string) Novel {
	bookId, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return Novel{
		TxtPath:  txtPath,
		Title:    title,
		Author:   author,
		Language: language,
		UUID:     bookId.String(),
	}
}

func editSuffix(novel Novel) {
	oldPath := novel.Title + ".zip"
	newPath := novel.Title + ".epub"
	err := os.Rename(oldPath, newPath)
	if err != nil {
		panic(err)
	}
}

func main() {

	txtPath := flag.String("path", "", "网文txt路径")
	title := flag.String("title", "", "网文名称")
	author := flag.String("author", "未知", "网文作者")

	flag.Parse()

	if *txtPath == "" || *title == "" {
		flag.Usage()
		return
	}

	novel := createNovel(*txtPath, *title, *author, "zh")

	// 删除生成文件夹
	clean(novel)

	// 创建书本文件夹以及子文件夹
	genBaseDir(novel)

	// 解析小说
	parseNovel(&novel)

	// 生成封面图
	genCover(novel)

	// 生成guide页面
	genGuide(novel)

	// 生成首页
	genInnerTitlePage(novel)

	// 生成样式文件
	genCss(novel)

	// 生成mimetype文件
	genMimetype(novel)

	// 生成container.xml文件
	genContainerXML(novel)

	// 生成content.opf文件
	genContentOpf(novel)

	// 生成toc.ncx文件
	genTocNcx(novel)

	// 压缩成zip文件
	genZip(novel)

	// 改后缀
	editSuffix(novel)

	// 生成azw3文件
	exportAzw3(novel)

	// 生成mobi文件
	exportMobi(novel)
}
