package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	ERR      = -1
	IsNoFile = 0
	IsAFile  = 1
	IsADir   = 2
)

const (
	RED    = "31"
	GREEN  = "32"
	YELLOW = "33"
	BLUE   = "34"
)

var COMMON = []string{"#", "//"}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		//  no such file or directory
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// 可能返回文件路径,可能返回目录路径,返回空代表文件和目录都不存在
func buildPath(path, keyword string) (stat int, filePath string) {

	filePath = Home + "/" + path + "/" + keyword + ".md"
	// 首先尝试寻找特定文件
	// 例如tellme -r vim, 则尝试在default(shell)目录下寻找vim.md
	// 例如tellme -r aws ec2, 则尝试在aws目录下寻找ec2.md
	// 例如tellme -r aws ec2 open, 则尝试在aws/ec2目录下寻找open.md
	if IsFile(filePath) {
		return IsAFile, filePath
	}

	// 在文件没找到的情况下path中包含"/" 说明有三个参数,属于精确查找
	// 例如tellme -r aws ec2 open
	if strings.Contains(path, "/") {
		return IsNoFile, filePath
	}

	// 如果没找到特定文件,则开始查找目录
	if path == DefaultPath {
		// 例如tellme -r vim, 则尝试在根目录下寻找vim文件夹
		filePath = Home + "/" + keyword
	} else {
		// 例如tellme -r aws ec2, 则尝试在根目录下寻找aws/ec2文件夹
		filePath = Home + "/" + path + "/" + keyword
	}

	if IsDir(filePath) {
		return IsADir, filePath
	}

	return IsNoFile, ""

}

func confirmInput(prompt string) bool {
	// 提示用户进行确认操作
	fmt.Print(prompt, "[yes/no]: ")

	// 读取用户输入
	var input string
	fmt.Scanln(&input)

	// 将输入转换为小写并去除首尾空格
	input = strings.TrimSpace(strings.ToLower(input))

	// 根据用户输入执行相应操作
	if input == "yes" {
		fmt.Println("操作已确认")
		return true
		// 在这里执行需要确认的操作
	} else if input == "no" {
		fmt.Println("操作已取消")
		return false
		// 在这里执行取消操作
	} else {
		fmt.Println("无效的输入")
		return false
	}
}

func markdownParse(markdown string, n int) (string, error) {

	if !strings.HasPrefix(markdown, DELIM) {
		return markdown, nil
	}

	parts := strings.Split(markdown, DELIM)

	if len(parts) < 4 {
		return markdown, fmt.Errorf("failed to delimit")
	}

	return parts[n], nil
}

func getTemplate() []byte {
	return []byte("---\ntag: [\"example_tag\"]\n---\n# <Title>\n---\n## Example")
}

func colorPrintln(color, text string) {
	fmt.Printf("\x1b[%sm%s\x1b[0m\n", color, text)
}

func colorPrint(text string) {

	for _, s := range COMMON {
		n := strings.Index(text, s)
		if n != -1 {
			fmt.Printf("\x1b[%sm%s\x1b[0m", YELLOW, text[:n])
			fmt.Printf("\x1b[%sm%s\x1b[0m", GREEN, text[n:])
			fmt.Println()
			return
		}
	}
	colorPrintln(YELLOW, text)
}

func colorfulPrint(text string) {

	for _, t := range strings.Split(text, "\n") {
		// 排除以```开头的行
		if strings.HasPrefix(t, "```") {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(t), "#") || strings.HasPrefix(strings.TrimSpace(t), "//") {
			colorPrintln(GREEN, t)
			continue
		}
		colorPrint(t)
	}
}

func columnPrint(list []string, operation func(string) string) {
	columnWidth := 20
	column := 5

	// 每行打印四个文件名
	for i := 0; i < len(list); i += column {
		for j := i; j < i+column && j < len(list); j++ {
			// 使用 Printf 格式化输出文件名，并保持列的宽度
			//fmt.Printf("%-*s", columnWidth, filepath.Base(files[j]))
			if operation == nil {
				fmt.Printf("%-*s", columnWidth, list[j])
			} else {
				fmt.Printf("%-*s", columnWidth, operation(list[j]))
			}
		}
		fmt.Println()
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func removeDuplication_map(arr []string) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}

func removeDuplication_sorted(arr []string) []string {
	length := len(arr)
	if length == 0 {
		return arr
	}

	j := 0
	for i := 1; i < length; i++ {
		if arr[i] != arr[j] {
			j++
			if j < i {
				swap(arr, i, j)
			}
		}
	}
	return arr[:j+1]
}

func swap(arr []string, a, b int) {
	arr[a], arr[b] = arr[b], arr[a]
}
