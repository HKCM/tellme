package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type TagStruct struct {
	Tag  string
	file FileS
}

type FileS struct {
	folder    string
	subfolder string
	filename  string
}

func updateDirTags(dirPath string) (err error) {

	d := getFileStruct(dirPath)
	keywords := make([]TagStruct, 0, 10)
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			slog.Debug("检查文件", "file", path)
			tags, err := getTagsFromFile(path)
			f := getFileStruct(path)
			if err != nil {
				fmt.Println(err)
			}
			for _, t := range tags {
				r := TagStruct{
					Tag:  t,
					file: f,
				}
				keywords = append(keywords, r)
			}
		}
		return err
	})

	writeTags(Home+"/"+d.folder+DBFILE, keywords)

	return nil
}

func writeTags(s string, keywords []TagStruct) {

	var f *os.File
	var err error
	// if IsFile(s) {
	// 	f, err = os.OpenFile(s, os.O_APPEND|os.O_WRONLY, 0644)
	// } else {
	// 	f, err = os.Create(s)
	// }
	f, err = os.Create(s)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()
	for _, k := range keywords {
		err = w.Write([]string{k.Tag, k.file.folder, k.file.subfolder, k.file.filename}) // calls Flush internally
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

func addIndex(index map[string][]string, tags []string, path string) map[string][]string {
	l := strings.LastIndex(path, "/")
	partDir := path[:l]

	if index == nil {
		index = make(map[string][]string)
	}

	for _, k := range tags {
		d, ok := index[k] // 查询是否存在关键词
		if !ok {
			index[k] = []string{path} // 如果不存在,则直接添加
		} else {
			index[k] = append(d, path)
		}
		sort.Strings(index[k])

		partPathKeyword := partDir + "/" + k // 设置目录关键词
		d, ok = index[partPathKeyword]
		if !ok {
			index[partPathKeyword] = []string{path}
		} else {
			index[partPathKeyword] = append(d, path)
		}
		sort.Strings(index[partPathKeyword])

	}

	d, ok := index[partDir] // 设置目录+关键词
	if !ok {
		index[partDir] = tags
	} else {
		index[partDir] = append(d, tags...)
	}
	sort.Strings(index[partDir])
	index[partDir] = removeDuplication_sorted(index[partDir])

	return index

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
