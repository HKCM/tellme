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

	slog.Debug("即将检查", "path", p.Path, "keyword", p.Keyword)
	stat, filePath := buildPath(p.Path, p.Keyword)
	// if stat == IsNoFile && p.ArgL ==1 {

	// }
	slog.Debug("即将进行分支判定", "filepath", filePath)
	switch stat {
	case IsAFile:
		cmdShowNote(filePath)
	case IsADir:
		slog.Debug("cmdShow 将显示目录下的文件")
		cmdShowFiles(filePath)
	case IsNoFile:
		slog.Debug("既不是文件也不是目录,将显示index文件")
		cmdShowIndex(p)
	}

}

func cmdShowIndex(p Parser) {

	var msg string
	indexKeyword := p.Path + "/" + p.Keyword

	slog.Debug("检查index.json文件", "indexKeyword", indexKeyword)

	index := readIndex()
	keywordPaths, find := index[indexKeyword]
	//find, keywordPaths := checkIndex(indexKeyword) // 如果没有相关文件和目录 则进行关键词检索

	if !find && p.ArgL == 1 {
		msg = fmt.Sprintf("索引中未发现关键词: %s", indexKeyword)
		slog.Debug("没有发现关键词,并且只有一个参数,将构建新的关键词", "关键词", indexKeyword, "新关键词", p.Keyword)
		indexKeyword = p.Keyword
		keywordPaths, find = index[indexKeyword]
	}

	if find {
		slog.Debug("发现关键字", "关键词", indexKeyword)
		if len(keywordPaths) == 1 { // 如果对应关键词下只有一条路径,则直接显示文件内容
			filePath := Home + "/" + keywordPaths[0]
			slog.Debug("关键字下只有一条记录,直接显示记录", "filePath", filePath)
			cmdShowNote(filePath)
			return
		}
		slog.Debug("关键字下存在多条记录,显示具体文件:\n", "关键词", indexKeyword)
		columnPrint(keywordPaths, nil)
		return
	}
	msg = fmt.Sprintf("索引中未发现关键词: %s", indexKeyword)

	slog.Debug(msg)
	slog.Warn(msg)

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

		key, _ := strings.CutSuffix(files[0], EXT)

		fmt.Printf("\n请使用类似命令获取详情: tellme %s %s\n", s, filepath.Base(key))
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

	text, part, err := markdownParse(string(markdown), 3)
	if err != nil {
		slog.Info("failed to parse markdown", "path", path, "err", err)
	}

	colorfulPrint(text)
	if part >= 4 {
		fmt.Println("Get more detail: ", path)
	}

}

func cmdEditNote(p Parser) {
	var cmd *exec.Cmd
	// 构建文件路径
	dirPath := Home + "/" + p.Path + "/" + p.Keyword
	filePath := Home + "/" + p.Path + "/" + p.Keyword + ".md"
	var oldTags []string
	var err error

	// 如果文件已存在直接编辑
	if IsFile(filePath) {
		slog.Debug("文件已存在", "文件路径", filePath)
		oldTags, err = getFileTags(filePath)
		if err != nil {
			slog.Warn("获取Tag失败", "filePath", filePath, "err", err)
		}
		cmd = exec.Command(EDITOR, filePath)
	} else {
		if IsDir(dirPath) {
			slog.Error("已存在同名目录", "目录", dirPath)
			return
		}
		fileDirPath := filepath.Dir(filePath)
		e, err := pathExists(fileDirPath)
		if err != nil {
			panic(err)
		}
		// 如果目录不存在则创建
		if !e {
			slog.Debug("目录不存在", "目录", fileDirPath)
			err = os.MkdirAll(fileDirPath, 0777)
			if err != nil {
				panic(err)
			}
		}
		var lineCmd string
		if IsFile(Template) {
			slog.Debug("模版文件已存在,使用模版...")
			lineCmd = fmt.Sprintf("autocmd VimEnter * nested silent! 0r %s", Template) // 使用template内容
		} else {
			slog.Debug("模版文件不存在,使用命令行参数添加模版...")
			lineCmd = "normal! i---\ntags: [\"example\"]\n---\n# <Title>\n---\n## Example"
		}
		slog.Debug(lineCmd)
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
		slog.Debug("如果文件不存在,则表示文件没有创建,不用更新索引")
		return
	}

	newTags, err := getFileTags(filePath)
	if err != nil {
		slog.Warn("新note获取Tag失败", "filePath", filePath, "err", err)
	}

	slog.Debug("新旧Tags:", "newTags", newTags, "oldTags", oldTags)
	if equal(newTags, oldTags) {
		slog.Debug("新旧文件Tag相同,不用更新索引")
		return
	}

	//TODO 更新索引
	slog.Info("即将更新索引:", "newTags", newTags, "oldTags", oldTags)
	index := readIndex()
	index = updateIndex(filePath, index, newTags, oldTags)
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
		if p.Confirm || confirmInput(fmt.Sprintf("即将删除文件 %s ", filePath)) {
			tags, err = getFileTags(filePath)
			if err != nil {
				panic(fmt.Errorf("获取文件Tag %s 失败,%v", filePath, err))
			}
			err = os.Remove(filePath)
			if err != nil {
				panic(fmt.Errorf("删除文件 %s 失败,%v", filePath, err))
			}
			indexUpdate = true
		}

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
