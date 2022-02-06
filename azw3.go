package main

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/leotaku/mobi"
	"golang.org/x/text/language"
)

func exportAzw3(novel Novel) {
	mb := mobi.Book{
		Title:       novel.Title,
		Authors:     []string{novel.Author},
		CreatedDate: time.Now(),
		Chapters:    []mobi.Chapter{},
		Language:    language.MustParse("zh"),
		UniqueID:    rand.Uint32(),
	}
	for _, chapter := range novel.Chapters {
		ch := mobi.Chapter{
			Title:  chapter.Title,
			Chunks: mobi.Chunks(chapter.XHTML),
		}
		mb.Chapters = append(mb.Chapters, ch)
	}

	mb.CSSFlows = []string{css}
	f, err := os.Open(filepath.Join(novel.Title, "images", "cover.jpg"))
	if err != nil {
		panic(fmt.Errorf("添加封面失败: %w", err))
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(fmt.Errorf("添加封面失败: %w", err))
	}
	mb.CoverImage = img

	// Convert book to PalmDB database
	db := mb.Realize()

	// Write database to file
	f, _ = os.Create(novel.Title + ".azw3")
	err = db.Write(f)
	if err != nil {
		panic(fmt.Errorf("保存失败: %w", err))
	}
}

func SectionSliceChunk(s []Chapter, size int) [][]Chapter {
	var ret [][]Chapter
	for size < len(s) {
		// s[:size:size] 表示 len 为 size，cap 也为 size，第二个冒号后的 size 表示 cap
		s, ret = s[size:], append(ret, s[:size:size])
	}
	ret = append(ret, s)
	return ret
}
