package main

import (
	"encoding/csv"
	"fmt"
	"io"
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

type TagData struct {
	tag       string
	folder    string
	subfolder string
	filename  string
}

func cmdShow(p Parser) {

	slog.Debug("即将检查", "path", p.Path, "tag", p.Tag)
	stat, filePath := buildPath(p.Path, p.Tag)
	// if stat == IsNoFile && p.ArgL ==1 {

	// }
	slog.Debug("即将进行分支判定", "filepath", filePath)
	switch stat {
	case IsAFile:
		cmdShowNote(filePath)
	case IsADir:
		if p.ShowTags {
			slog.Debug("cmdShow 将显示目录下的Tags")
			cmdShowDirTags(filePath)
		} else {
			slog.Debug("cmdShow 将显示目录下的文件")
			cmdShowFiles(filePath)
		}

	case IsNoFile:
		slog.Debug("既不是文件也不是目录,将显示tags文件")
		cmdShowTag(p)
	}

}

func cmdShowTag(p Parser) {

	tagFiles := getTagsFromDir(Home + "/" + p.Path)

	var filename []string

	for _, tf := range tagFiles {
		if tf.tag == p.Tag {
			filename = append(filename, tf.filename)
		}
	}

	if len(filename) == 0 {
		msg := fmt.Sprintf("索引中未发现关键词: %s", p.Tag)
		slog.Warn(msg)
		return
	}

	if len(filename) == 1 {
		cmdShowNote(Home + "/" + tagFiles[0].folder + "/" + tagFiles[0].subfolder + "/" + tagFiles[0].filename)
		return
	}

	if len(filename) > 1 {
		columnPrint(filename, nil)
		return
	}
}

func cmdShowDirTags(path string) error {
	tagFiles := getTagsFromDir(path)

	d := getDirStruct(path)

	if len(tagFiles) == 0 {
		slog.Error("There are no db file")
		return nil
	}

	var tags []string

	for _, tf := range tagFiles {
		tags = append(tags, tf.tag)
	}
	fmt.Printf("存在以下Tags:\n")
	columnPrint(tags, nil)

	fmt.Printf("\n请使用类似命令获取详情: tellme %s %s\n\n", d.folder, tags[0])

	return nil
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

		fmt.Printf("\n请使用类似命令获取详情: tellme %s %s\n\n", s, filepath.Base(key))
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
	dirPath := Home + "/" + p.Path + "/" + p.Tag
	filePath := Home + "/" + p.Path + "/" + p.Tag + ".md"
	var allTagDatas []TagData
	var oldTags []string
	// var oldTagDatas []TagData
	var err error

	// 如果文件已存在直接编辑
	if IsFile(filePath) {
		slog.Debug("文件已存在", "文件路径", filePath)

		allTagDatas = getAllTagData(filePath)
		for _, t := range allTagDatas {
			if t.filename == p.Tag+".md" {
				oldTags = append(oldTags, t.tag)
				// oldTagDatas = append(oldTagDatas, t)
			}
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

	newTags, err := getTagsFromFile(filePath)
	if err != nil {
		slog.Warn("新note获取Tag失败", "filePath", filePath, "err", err)
	}

	slog.Debug("新旧Tags:", "newTags", newTags, "oldTags", oldTags)
	if equal(newTags, oldTags) {
		slog.Debug("新旧文件Tag相同,不用更新索引")
		return
	}

	//TODO 更新索引
	f := getFileStruct(filePath)

	var newTagDatas []TagData
	for _, tag := range newTags {
		newTagDatas = append(newTagDatas, TagData{tag: tag, folder: f.folder, subfolder: f.subfolder, filename: f.filename})
	}

	if len(oldTags) == 0 {
		appendTags(Home+"/"+f.folder+DBFILE, newTagDatas)
	} else {
		for _, t := range allTagDatas {
			if t.filename != f.filename {
				newTagDatas = append(newTagDatas, t)
			}
		}
		writeTags(Home+"/"+f.folder+DBFILE, newTagDatas)
	}

}

func appendTags(dbPath string, ts []TagData) {

	var f *os.File
	var err error
	if IsFile(dbPath) {
		f, err = os.OpenFile(dbPath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		f, err = os.Create(dbPath)
	}
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	defer w.Flush()
	for _, k := range ts {
		err = w.Write([]string{k.tag, k.folder, k.subfolder, k.filename}) // calls Flush internally
		if err != nil {
			panic(err)
		}
	}
}

func getTagsFromDir(dirPath string) (tagFiles []TagData) {
	// 获取文件内容
	d := getDirStruct(dirPath)
	fdb, err := os.Open(Home + "/" + d.folder + DBFILE)

	if err != nil {
		panic(err)
	}

	r := csv.NewReader(fdb)
	for {

		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		if d.subfolder != "" {
			if d.subfolder == record[2] {
				tagFiles = append(tagFiles, TagData{tag: record[0], folder: record[1], subfolder: record[2], filename: record[3]})
			}
		} else {
			tagFiles = append(tagFiles, TagData{tag: record[0], folder: record[1], subfolder: record[2], filename: record[3]})
		}

	}

	return
}

// 从数据库中获取Note的Tags
// func getTagsByNote(filePath string) (tags []string) {
// 	// 获取文件内容
// 	f := getFileStruct(filePath)
// 	fdb, err := os.Open(Home + "/" + f.folder + DBFILE)
// 	if err != nil {
// 		panic(err)
// 	}
// 	r := csv.NewReader(fdb)
// 	for {
// 		record, err := r.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			panic(err)
// 		}
// 		if record[3] == filepath.Base(filePath) {
// 			tags = append(tags, record[0])
// 		}
// 	}
// 	return
// }

func getAllTagData(filePath string) (tagDatas []TagData) {
	f := getFileStruct(filePath)
	dbFile := Home + "/" + f.folder + DBFILE
	if !IsFile(dbFile) {
		slog.Debug("数据库文件不存在", "文件路径", dbFile)
		return
	}
	fdb, err := os.Open(dbFile)
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(fdb)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		tagDatas = append(tagDatas, TagData{tag: record[0], folder: record[1], subfolder: record[2], filename: record[3]})
	}
	return
}

// 删除必须要指定note的名字
func cmdRemoveNote(p Parser) {

	stat, filePath := buildPath(p.Path, p.Tag)
	var err error
	switch stat {
	case IsAFile:
		if p.Confirm || confirmInput(fmt.Sprintf("即将删除文件 %s ", filePath)) {
			err = os.Remove(filePath)
			if err != nil {
				panic(fmt.Errorf("删除文件 %s 失败,%v", filePath, err))
			}
			f := getFileStruct(filePath)
			var newTagDatas []TagData
			tagDatas := getAllTagData(Home + "/" + f.folder + DBFILE)
			for _, t := range tagDatas {
				if t.filename == f.filename && t.subfolder == f.subfolder {
					continue
				}
				newTagDatas = append(newTagDatas, t)
			}
			writeTags(Home+"/"+f.folder+DBFILE, newTagDatas)
		}

	case IsADir:
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
			err = os.RemoveAll(filePath)
			if err != nil {
				panic(err)
			}
			d := getDirStruct(filePath)
			if d.subfolder != "" {
				var newTagDatas []TagData
				tagDatas := getAllTagData(Home + "/" + d.folder + DBFILE)
				for _, t := range tagDatas {
					if t.subfolder == d.subfolder {
						continue
					}
					newTagDatas = append(newTagDatas, t)
				}
				writeTags(Home+"/"+d.folder+DBFILE, newTagDatas)
			}
		}
	case IsNoFile:
		fmt.Printf("没有发现相关文件或目录")
	default:
		panic(fmt.Errorf("未知状态"))
	}

}

func cmdInit(p Parser) {
	// 创建template文件
	createTemplate(p.Confirm)
}

func createTemplate(confirmed bool) {
	if !IsFile(Template) || confirmed || confirmInput(fmt.Sprintf("模版文件 %s 已存在, 确认使用初始化覆盖?", Template)) {
		err := os.WriteFile(Template, getTemplate(), 0666) //写入文件(字节数组)
		if err != nil {
			panic(err)
		}
		return
	}
}

func cmdShowHelp() {
	fmt.Println(usage())
	os.Exit(0)
}

func cmdTagsUpdate(p Parser) {

	dirPath := Home + "/"

	if !IsDir(dirPath) {
		slog.Error("未发现待更新的目录", "目录", dirPath)
		return
	}
	slog.Debug("即将更新", "目录", dirPath)
	if strings.ToLower(p.UpdateTarget) == "all" {
		fs, err := os.ReadDir(dirPath)
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			if f.IsDir() {
				updateDirTags(Home + "/" + f.Name())
			}
		}
	}

	if strings.ToLower(p.UpdateTarget) != "all" {
		dirPath = Home + "/" + p.UpdateTarget
		updateDirTags(dirPath)
	}

}
