package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func checkIndex(keyword string) (find bool, dirPaths []string) {

	jsonContent, err := os.ReadFile(IndexFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	// 定义一个 map 用于存储解析后的 JSON 数据
	data := make(map[string][]string)

	// 将 JSON 字符串解析到 map 中
	err = json.Unmarshal(jsonContent, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	dirPaths, find = data[keyword]
	return
}

func updateIndex(filePath string, newTags, oldTags []string) (data map[string][]string) {

	partFilePath, _ := strings.CutPrefix(filePath, home+"/")

	data = readIndex()

	data = deleteIndex(data, oldTags, partFilePath)
	data = addIndex(data, newTags, partFilePath)

	return data
}

func readIndex() (data map[string][]string) {
	var err error
	var jsonContent []byte
	if !IsFile(IndexFile) { //如果文件不存在
		return
	}

	jsonContent, err = os.ReadFile(IndexFile)
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(IndexFile, os.O_WRONLY|os.O_TRUNC, 0666) //打开文件
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 定义一个 map 用于存储解析后的 JSON 数据

	// 将 JSON 字符串解析到 map 中
	err = json.Unmarshal(jsonContent, &data)
	if err != nil {
		panic(err)
	}

	return
}

func writeIndex(data map[string][]string) {
	var file *os.File
	var err error
	if !IsFile(IndexFile) { //如果文件不存在
		file, err = os.Create(IndexFile) //创建文件
		fmt.Printf("%s 文件不存在,创建文件并写入\n", IndexFile)
	} else {
		file, err = os.OpenFile(IndexFile, os.O_WRONLY|os.O_TRUNC, 0666) //打开文件
		fmt.Printf("%s 文件存在,写入文件\n", IndexFile)
	}
	if err != nil {
		panic(err)
	}

	content, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Errorf("将map转化为content出错 %v", err))
	}
	writer := bufio.NewWriter(file)

	_, err = writer.Write(content)
	if err != nil {
		panic(fmt.Errorf("index 写入文件出错 %v", err))
	}
	writer.Flush()
}

func updateAllIndex() error {
	// 获取新的文件内容
	return filepath.WalkDir(IndexFile, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			fmt.Printf("%s 是一个目录\n", path)
		} else {
			fmt.Printf("%s 是一个文件\n", path)
		}
		return err
	})
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

func deleteIndex(index map[string][]string, tags []string, path string) map[string][]string {
	l := strings.LastIndex(path, "/")
	partDir := path[:l]

	for _, k := range tags {
		// 关键词移除
		{
			d, ok := index[k]
			if ok {
				if len(d) == 1 { // 如果有且只有一个关键词 则直接删除
					delete(index, k)
				} else {
					var arr []string
					for _, i := range d {
						if i != path {
							arr = append(arr, i)
						}
					}
					index[k] = arr
				}
			}
		}

		// 移除目录关键词
		{
			partPathKeyword := partDir + "/" + k
			d, ok := index[partPathKeyword]
			if ok {
				if len(d) == 1 {
					delete(index, partPathKeyword)
				} else {
					var arr []string
					for _, i := range d {
						if i != path {
							arr = append(arr, i)
						}
					}
					index[partPathKeyword] = arr
				}
			}
		}

		// 移除目录+关键词
		{
			d, ok := index[partDir]
			if ok {
				if len(d) == 1 {
					delete(index, partDir)
				} else {
					var arr []string
					for _, i := range d {
						if i != k {
							arr = append(arr, i)
						}
					}
					index[partDir] = arr
				}
			}
		}

	}

	return index

}

func getFileTags(filePath string) (tags []string, err error) {
	// 获取文件内容
	markdown, err := os.ReadFile(filePath)
	if err != nil {
		err = fmt.Errorf("failed to read file: %s, %v", filePath, err)
		return
	}

	text, err := markdownParse(string(markdown), 1) // 获取tag
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
	fmt.Println(tag.Tags)
	tags = tag.Tags
	return

}
