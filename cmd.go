package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CommandType int

const (
	NoCmd = iota
	HelpCmd
	EditCmd
	ShowCmd
	RmCmd
)

func selectCmd(model Model) (CommandType, error) {

	switch model {
	case HelpModel:
		return HelpCmd, nil
	case RemoveModel:
		return RmCmd, nil
	case EditModel:
		return EditCmd, nil
	case NormalModel:
		return ShowCmd, nil
	default:
		fmt.Println("Model is:", model)
		return NoCmd, errors.New("UnknownModel")
	}

}

func checkPath(path, keyword string) (stat int, filePath, dirPath string) {
	// 首先尝试寻找特定文件
	filePath = home + "/" + path + "/" + keyword + ".md"

	stat, err := getPathStat(filePath)
	if err != nil {
		panic(err)
	}

	if stat == IsAFile {
		return IsAFile, filePath, path
	}

	// 如果没找到特定文件,则查找目
	if path == "" {
		dirPath = home + "/" + keyword
	} else {
		dirPath = home + "/" + path + "/" + keyword
	}

	stat, err = getPathStat(dirPath)
	if err != nil {
		panic(err)
	}
	return stat, filePath, dirPath

}

func cmdShow(path, keyword string) {
	stat, filePath, dirPath := checkPath(path, keyword)

	switch stat {
	case IsAFile:
		cmdShowNote(filePath)
	case IsADir:
		cmdShowList(dirPath)
	default:
		fmt.Printf("没有发现 %s 文件或 %s 目录\n", filePath, dirPath)
	}
}
func cmdShowList(path string) {
	files, err := filepath.Glob(path + "/*")
	if err != nil {
		panic(err)
	}
	if len(files) > 0 {
		fmt.Println("目录下存在以下文件:")
		for _, f := range files {
			fmt.Println(f)
		}

		part := path[len(home):]
		s := strings.ReplaceAll(part, "/", " ")
		fmt.Printf("请使用以下命令: tellme %s <Name>", s)
	} else {
		fmt.Printf("目录 %s 下为空\n", path)
	}
}

func cmdShowNote(path string) {

	markdown, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("failed to read file: %s, %v", path, err)
		return
	}

	text, err := mdParse(string(markdown))
	if err != nil {
		fmt.Printf("failed to parse markdown: %s, %v", path, err)
	}
	fmt.Println(text)
}

func mdParse(markdown string) (string, error) {
	delim := "---\n"

	// if the markdown does not contain frontmatter, pass it through unmodified
	if !strings.HasPrefix(markdown, delim) {
		return markdown, nil
	}

	// otherwise, split the frontmatter and cheatsheet text
	parts := strings.Split(markdown, delim)

	if len(parts) < 4 {
		return markdown, fmt.Errorf("failed to delimit")
	}

	//fmt.Println(parts)

	return parts[3], nil
}

func cmdEditNote(path, keyword string) {

	if path == "" {
		fmt.Println("当使用编辑模式时,至少包含两个层级: tellme -e vim open")
		return
	}

	stat, filePath, dirPath := checkPath(path, keyword)
	var cmd *exec.Cmd
	switch stat {
	case IsAFile:
		cmd = exec.Command(editor, filePath)
	case IsADir:
		fmt.Printf("已存在同名目录: %s\n", dirPath)
		fmt.Printf("Note创建失败: %s\n", filePath)
		return
	case IsNoFile:
		e, err := pathExists(filepath.Dir(filePath))
		if err != nil {
			panic(err)
		}
		// 如果目录不存在则创建
		if !e {
			os.MkdirAll(dirPath, 0777)
		}
		lineCmd := fmt.Sprintf("normal! i---\ntag: [\"%s\",\"example\"]\n---\n# <Title>\n---\n## Example", keyword)
		//fmt.Println(lineCmd)
		cmd = exec.Command(editor, "-c", lineCmd, filePath)
	default:
		fmt.Println("出现未知类型...")
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to edit %s: %v\n", filePath, err)
		os.Exit(1)
	}
	//TODO 更新索引

}

// 删除必须要指定note的名字
func cmdRemoveNote(path, keyword string) {

	stat, filePath, dirPath := checkPath(path, keyword)

	switch stat {
	case IsAFile:
		if confirm(fmt.Sprintf("即将删除 %s 文件操作？(yes/no): ", filePath)) {
			os.Remove(filePath)
		}
		//TODO 更新索引
		return
	case IsADir:
		files, err := filepath.Glob(dirPath + "/*")
		if err != nil {
			panic(err)
		}
		if len(files) > 0 {
			fmt.Println("目录下存在以下文件:")
			for _, f := range files {
				fmt.Println(f)
			}
		}
		if confirm(fmt.Sprintf("即将删除 %s 目录操作？(yes/no): ", dirPath)) {
			err = os.RemoveAll(dirPath)
			if err != nil {
				panic(err)
			}
		}
	case IsNoFile:
		fmt.Printf("没有发现 %s 文件或 %s 目录", filePath, dirPath)
	}
}

func cmdShowHelp(path, keyword string) {
	usage()
}

func confirm(prompt string) bool {
	// 提示用户进行确认操作
	fmt.Print(prompt)

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
