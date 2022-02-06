package main

import (
	"path/filepath"

	"github.com/766b/mobi"
)

func exportMobi(novel Novel) {
	m, err := mobi.NewWriter(novel.Title + ".mobi")
	if err != nil {
		panic(err)
	}

	m.Title(novel.Title)
	m.Compression(mobi.CompressionNone)

	coverPath := filepath.Join(novel.Title, "images", "cover.jpg")
	m.AddCover(coverPath, coverPath)

	m.NewExthRecord(mobi.EXTH_DOCTYPE, "EBOK")
	m.NewExthRecord(mobi.EXTH_AUTHOR, novel.Author)

	for _, chapter := range novel.Chapters {
		m.NewChapter(chapter.Title, []byte(chapter.XHTML))
	}

	m.Write()
}
