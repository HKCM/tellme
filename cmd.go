package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CommandType int
type ShowType int

const (
	InitCmd CommandType = iota
	HelpCmd
	EditCmd
	ShowCmd
	RmCmd
)
const (
	ShowNote ShowType = iota
	ShowIndex
	ShowFiles
	ShowNothing
)

func cmdShow(p Parser) {

	stat, filePath := buildPath(p.Path, p.Keyword)
	slog.Debug("即将检查", filePath)
	switch stat {
	case IsAFile:
		cmdShowNote(filePath)
		fmt.Println("for more detail:", filePath)
	case IsADir:
		slog.Debug("cmdShow 将显示目录下的文件")
		cmdShowFiles(filePath)
	case IsNoFile:
		slog.Debug("cmdShow 将显示index文件")
		cmdShowIndex(p.Path, p.Keyword)
	}

}

func cmdShowIndex(path, keyword string) {
	indexKeyword := path + "/" + keyword
	find, keywordPaths := checkIndex(indexKeyword) // 如果没有相关文件和目录 则进行关键词检索

	if find {
		if len(keywordPaths) == 1 { // 如果对应关键词下只有一条路径,则直接显示文件内容
			realPath := Home + "/" + keywordPaths[0]
			cmdShowNote(realPath)
			return
		}
		slog.Debug("发现关键字,属于以下目录:\n", "关键词", indexKeyword)
		columnPrint(keywordPaths, nil)
		return
	}

	slog.Debug("没有发现关键字,什么都不显示", "关键词", indexKeyword)
	cmdShowNothing(path, keyword)

}

func cmdShowNothing(path, keyword string) {
	fmt.Printf("未发现 %s 目录...\n", getRelativePath(path))
	fmt.Printf("搜索完毕,未发现任何 %s 相关记录...\n", keyword)
}

func getRelativePath(absPath string) string {
	return strings.Replace(absPath, Home, ShortHome, 1)
}

func cmdShowFiles(path string) error {
	files, err := filepath.Glob(path + "/*")
	if err != nil {
		return err
	}
	if len(files) > 0 {
		part := path[len(Home)+1:] // 加1去除末尾的 “/”
		s := strings.ReplaceAll(part, "/", " ")
		fmt.Printf("%s 目录下存在以下文件(%d):\n", path, len(files))

		columnPrint(files, filepath.Base)

		fmt.Printf("\n请使用类似命令获取详情: tellme %s %s\n", s, filepath.Base(files[0][:len(files[0])-3]))
	} else {
		fmt.Printf("目录 %s 下为空\n", path)
	}

	return nil
}

func cmdShowNote(path string) {

	markdown, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %s, %v", path, err))
	}

	text, err := markdownParse(string(markdown), 3)
	if err != nil {
		slog.Info("failed to parse markdown", "path", path, "err", err)
	}

	colorfulPrint(text)

}

func cmdEditNote(p Parser) {
	var cmd *exec.Cmd
	// 构建文件路径
	filePath := Home + "/" + p.Path + "/" + p.Keyword + ".md"
	var originTags []string
	var err error

	// 如果文件已存在直接编辑
	if IsFile(filePath) {
		fmt.Printf("文件 %s 已存在\n", filePath)
		originTags, err = getFileTags(filePath)
		if err != nil {
			panic(fmt.Errorf("删除文件 %s 失败,%v", filePath, err))
		}
		cmd = exec.Command(EDITOR, filePath)
	} else {
		dirPath := filepath.Dir(filePath)
		e, err := pathExists(dirPath)
		if err != nil {
			panic(err)
		}
		// 如果目录不存在则创建
		if !e {
			fmt.Println("目录:", dirPath, "不存在,创建中...")
			err = os.MkdirAll(dirPath, 0777)
			if err != nil {
				panic(err)
			}
		}
		var lineCmd string
		if IsFile(Template) {
			fmt.Println("模版文件已存在,使用模版...")
			lineCmd = fmt.Sprintf("autocmd VimEnter * nested silent! 0r %s", Template) // 使用template内容
		} else {
			fmt.Println("模版文件不存在,使用命令行参数添加模版...")
			lineCmd = "normal! i---\ntags: [\"example\"]\n---\n# <Title>\n---\n## Example"
		}
		fmt.Println(lineCmd)
		cmd = exec.Command(EDITOR, "-c", lineCmd, filePath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("failed to edit %s: %v", filePath, err))
	}

	// 如果文件不存在,则表示文件没有创建,不用更新索引
	if !IsFile(filePath) {
		fmt.Println("如果文件不存在,则表示文件没有创建,不用更新索引")
		return
	}

	newTags, err := getFileTags(filePath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("newTags:%v oldTags:%v \n", newTags, originTags)
	if equal(newTags, originTags) {
		fmt.Println("新旧文件Tag相同,不用更新索引")
		return
	}

	//TODO 更新索引
	index := readIndex()
	index = updateIndex(filePath, index, newTags, originTags)
	writeIndex(index)

}

// 删除必须要指定note的名字
func cmdRemoveNote(p Parser) {

	indexUpdate := false
	stat, filePath := buildPath(p.Path, p.Keyword)
	var tags []string
	var err error

	switch stat {
	case IsAFile:
		if p.Confirm || confirmInput(deleteFileMsg(p, filePath)) {
			tags, err = getFileTags(filePath)
			if err != nil {
				panic(fmt.Errorf("获取文件Tag %s 失败,%v", filePath, err))
			}
			err = os.Remove(filePath)
			if err != nil {
				panic(fmt.Errorf("删除文件 %s 失败,%v", filePath, err))
			}
		}
		indexUpdate = true

	case IsADir: // TODO 目前删除目录有问题
		files, err := filepath.Glob(filePath + "/*")
		if err != nil {
			panic(err)
		}
		if len(files) > 0 {
			fmt.Println("目录下存在以下文件:")
			for _, f := range files {
				fmt.Println(f)
			}
		}
		if confirmInput(fmt.Sprintf("即将删除 %s 目录操作", filePath)) {
			updateIndexByRemoveDir(filePath)
			err = os.RemoveAll(filePath)
			if err != nil {
				panic(err)
			}
		}
		indexUpdate = true
	case IsNoFile:
		fmt.Printf("没有发现相关文件或目录")
	default:
		panic(fmt.Errorf("未知状态"))
	}

	if indexUpdate {
		//TODO 更新索引
		index := readIndex()
		data := updateIndex(filePath, index, []string{}, tags)
		writeIndex(data)
	}

}

// 创建template文件
func cmdInit(p Parser) {
	if !IsFile(Template) || p.Confirm || confirmInput(fmt.Sprintf("模版文件 %s 已存在, 确认使用初始化覆盖?", Template)) {
		err := os.WriteFile(Template, getTemplate(), 0666) //写入文件(字节数组)
		if err != nil {
			panic(err)
		}
		return
	}

}

func cmdShowHelp(p Parser) {
	fmt.Println(usage())
	os.Exit(1)
}

func cmdIndexUpdate(p Parser) {
	updateAllIndex()

}
