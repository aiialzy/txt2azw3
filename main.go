package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/golang/freetype"
	"github.com/google/uuid"
)

type Novel struct {
	UUID     string
	Title    string
	Author   string
	Language string
	Chapters []Chapter
}

func genCover(novel Novel) {
	// 图片宽度
	width := 600

	// 图片高度
	height := 800

	imgFile, _ := os.Create(path.Join(novel.Title, "images", "cover.jpg"))
	defer imgFile.Close()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 设置背景色
	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			img.Set(y, x, color.RGBA{255, 255, 255, 255})
		}
	}

	fontBytes, err := ioutil.ReadFile("font.ttf")
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	f := freetype.NewContext()

	f.SetDPI(300)

	f.SetFont(font)

	f.SetFontSize(26)

	f.SetClip(img.Bounds())

	f.SetDst(img)

	f.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 255}))

	pt := freetype.Pt(20, 400)

	_, err = f.DrawString(novel.Title, pt)
	if err != nil {
		panic(err)
	}

	err = jpeg.Encode(imgFile, img, nil)
	if err != nil {
		panic(err)
	}
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

func createNovel(title, author, language string) Novel {
	bookId, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return Novel{
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

	novel := createNovel("凡人修仙传", "忘语", "zh")

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
