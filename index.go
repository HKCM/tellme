package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type FileS struct {
	folder    string
	subfolder string
	filename  string
}

func updateDirTags(dirPath string) (err error) {

	d := getFileStruct(dirPath)
	keywords := make([]TagData, 0, 10)
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			slog.Debug("检查文件", "file", path)
			tags, err := getTagsFromFile(path)
			f := getFileStruct(path)
			if err != nil {
				fmt.Println(err)
			}
			for _, t := range tags {
				r := TagData{
					tag:       t,
					folder:    f.folder,
					subfolder: f.subfolder,
					filename:  f.filename,
				}
				keywords = append(keywords, r)
			}
		}
		return err
	})

	writeTags(Home+"/"+d.folder+DBFILE, keywords)

	return nil
}

func writeTags(dbFile string, tagDatas []TagData) {

	var f *os.File
	var err error
	f, err = os.Create(dbFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()
	for _, k := range tagDatas {
		err = w.Write([]string{k.tag, k.folder, k.subfolder, k.filename}) // calls Flush internally
		if err != nil {
			panic(err)
		}
	}
}

func getDirStruct(path string) FileS {
	p, _ := strings.CutPrefix(path, Home+"/")
	s := strings.Split(p, "/")
	n := len(s)
	var f FileS
	f.folder = s[0]
	if n == 2 {
		f.subfolder = s[1]
	}
	return f
}

func getFileStruct(path string) FileS {
	p, _ := strings.CutPrefix(path, Home+"/")
	s := strings.Split(p, "/")
	n := len(s)
	var f FileS
	f.folder = s[0]
	f.filename = s[n-1]
	if n == 3 {
		f.subfolder = s[1]
	}
	return f
}

func getTagsFromFile(filePath string) (tags []string, err error) {
	// 获取文件内容
	markdown, err := os.ReadFile(filePath)
	if err != nil {
		err = fmt.Errorf("failed to read file: %s, %v", filePath, err)
		return
	}

	text, _, err := markdownParse(string(markdown), 1) // 获取tag
	if err != nil {
		err = fmt.Errorf("failed to parse markdown: %s, %v, got the hole file", filePath, err)
		return
	}

	tag := struct {
		Tags []string `yaml:"tags"`
	}{}

	if err = yaml.Unmarshal([]byte(text), &tag); err != nil {
		return
	}
	slog.Debug("获取到文件tag", "文件", filePath, "tags", tag.Tags)
	tags = tag.Tags
	return

}
