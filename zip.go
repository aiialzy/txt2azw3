package main

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func genZip(novel Novel) {
	zipFile, err := os.Create(novel.Title + ".zip")
	if err != nil {
		panic(err)
	}

	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(novel.Title, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == novel.Title {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, novel.Title+"\\")
		header.Name = strings.TrimPrefix(header.Name, novel.Title+"/")

		if info.IsDir() {
			header.Name += "/"
		} else {
			if header.Name == "mimetype" {
				header.Method = zip.Store
			} else {
				header.Method = zip.Deflate
			}
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return err
	})

	if err != nil {
		panic(err)
	}
}
