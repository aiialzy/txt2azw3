package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path"

	"github.com/golang/freetype"
)

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

	fontBytes, err := root.ReadFile("font.ttf")
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
