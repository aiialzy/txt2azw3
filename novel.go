package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"text/template"

	cn "github.com/aiialzy/chinese-number"
)

type Chapter struct {
	Index    int
	Title    string
	Content  string
	TextPath string
	Filename string
	XHTML    string
}

func parseNovel(novel *Novel) {
	txtBytes, err := os.ReadFile(novel.TxtPath)
	if err != nil {
		panic(err)
	}

	txt := string(txtBytes)

	txt = strings.ToValidUTF8(txt, "")

	numPart := "[零一二三四五六七八九十百千万亿壹貳叁肆伍陸柒捌玖拾佰仟白干〇兩两0-9a-zA-Z]+章"

	reg := regexp.MustCompile(`第.*?` + numPart + `[^，。,.]*?(\n|$|\r)`)
	numReg := regexp.MustCompile(numPart)
	titleReg := regexp.MustCompile(`[.<>?|*:"]`)

	resultIndex := reg.FindAllStringIndex(txt, -1)
	lastBegin := 0

	for i, index := range resultIndex {
		line := txt[index[0]:index[1]]
		chapterTitle := strings.TrimFunc(line, func(ch rune) bool {
			return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\f' || ch == '\r'
		})

		found := numReg.FindAllStringSubmatch(chapterTitle, -1)
		numRunes := []rune(found[0][0])
		numRunes = numRunes[:len(numRunes)-1]
		chapterIndex, err := cn.Parse(string(numRunes))

		if err != nil {
			panic(err)
		}

		chapterTitle = titleReg.ReplaceAllString(chapterTitle, " ")

		filename := fmt.Sprintf("chapter%d", chapterIndex)
		spaceReg := regexp.MustCompile("[\t \r\n\f]")
		filename = spaceReg.ReplaceAllString(filename, "_")
		filename += ".xhtml"

		textPath := path.Join(novel.Title, "text", filename)

		chapter := Chapter{
			Index:    int(chapterIndex),
			Title:    chapterTitle,
			TextPath: textPath,
			Filename: filename,
		}

		if i > 0 {
			novel.Chapters[i-1].Content = txt[lastBegin:index[0]]
			lastBegin = index[1]
		}

		novel.Chapters = append(novel.Chapters, chapter)
	}

	novel.Chapters[len(novel.Chapters)-1].Content = txt[lastBegin:]

	novel.Chapters = removeDuplicate(novel.Chapters)

	sort.Slice(novel.Chapters, func(i, j int) bool {
		return novel.Chapters[i].Index < novel.Chapters[j].Index
	})

	for i, chapter := range novel.Chapters {
		novel.Chapters[i].XHTML = genXHTML(novel.Title, chapter)
	}

}

func removeDuplicate(novel []Chapter) (result []Chapter) {
	m := make(map[int]Chapter)
	for _, chapter := range novel {
		inmapChapter, ok := m[chapter.Index]
		if !ok || len(chapter.Content) > len(inmapChapter.Content) {
			m[chapter.Index] = chapter
		}
	}

	for _, chapter := range m {
		result = append(result, chapter)
	}

	return
}

func genXHTML(title string, chapter Chapter) string {
	t, err := template.ParseFiles("templates/text.tmpl")
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(chapter.TextPath, os.O_CREATE|os.O_WRONLY, 0766)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	lines := strings.Split(chapter.Content, "\r\n")
	var notEmptyLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			notEmptyLines = append(notEmptyLines, line)
		}
	}

	args := map[string]interface{}{
		"title":        title,
		"chapterId":    fmt.Sprintf("chapter%d", chapter.Index),
		"chapterTitle": chapter.Title,
		"sections":     notEmptyLines,
	}

	var xhtmlBuf bytes.Buffer

	err = t.Execute(io.MultiWriter(&xhtmlBuf, f), args)
	if err != nil {
		panic(err)
	}

	return xhtmlBuf.String()

}
