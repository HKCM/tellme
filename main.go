package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/user"
)

const (
	EDITOR = "vim"
	DELIM  = "---\n"
)

var (
	Home      string // u.HomeDir + /.tellme
	Template  string
	IndexFile string
	i         input
)

var (
	helpModel    bool
	initModel    bool
	indexModel   bool
	debugModel   bool
	editModel    bool
	removeModel  bool
	confirmModel bool
)

const Deep = 3

const (
	DefaultPath   = "shell" // 默认的查询路径
	Root          = "/.tellme"
	ShortHome     = "~" + Root
	ShortIndex    = ShortHome + "/index.json"
	ShortTemplate = ShortHome + "/template"
	EXT           = ".md"
)

func init() {

	u, err := user.Current()
	if err != nil {
		slog.Error("获取用户失败", err)
		os.Exit(1)
	}
	Home = u.HomeDir + Root
	if !IsDir(Home) {
		slog.Error(fmt.Sprintf("项目 %s 目录不存在,请参考README.md", Home))
		os.Exit(1)
	}
	Template = Home + "/template"
	IndexFile = Home + "/index.json"

	flag.BoolVar(&helpModel, "h", false, "显示帮助")
	flag.BoolVar(&initModel, "init", false, "创建模版")
	flag.BoolVar(&indexModel, "make-index", false, "更新索引")
	flag.BoolVar(&debugModel, "debug", false, "排错模式")
	flag.BoolVar(&editModel, "e", false, "编辑或新建笔记")
	flag.BoolVar(&removeModel, "r", false, "移除笔记")
	flag.BoolVar(&confirmModel, "y", false, "确认操作")
	flag.Parse()
	if debugModel {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	i = input{
		helpModel:   helpModel,
		initModel:   initModel,
		indexModel:  indexModel,
		editModel:   editModel,
		removeModel: removeModel,
		confirm:     confirmModel,
		args:        flag.Args(),
	}

	slog.Debug("input", "i", i)
}

func main() {

	p, err := parser(i)
	if err != nil {
		slog.Debug("parser 有错误, 设置为help模式%v", err)
		p.Model = ModelHelp
	}

	var runCmd func(p Parser)

	switch p.Model {
	case ModelInit:
		runCmd = cmdInit
	case ModelIndex:
		runCmd = cmdIndexUpdate
	case ModelHelp:
		runCmd = cmdShowHelp
	case ModelEdit:
		runCmd = cmdEditNote
	case ModelRemove:
		runCmd = cmdRemoveNote
	default:
		runCmd = cmdShow
	}
	runCmd(p)
}
