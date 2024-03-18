package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
)

const (
	editor = "vim"
	delim  = "---\n"
)

var (
	home         string // u.HomeDir + /.tellme
	realTemplate string
	IndexFile    string
	i            input
)

var (
	helpModel    bool
	initModel    bool
	indexModel   bool
	editModel    bool
	removeModel  bool
	confirmModel bool
)

const Deep = 3

const (
	defaultPath   = "shell" // 默认的查询路径
	root          = "/.tellme"
	shortHome     = "~" + root
	shortIndex    = shortHome + "/index.json"
	shortTemplate = shortHome + "/template"
)

func init() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	home = u.HomeDir + root
	realTemplate = home + "/template"
	IndexFile = home + "/index.json"

	flag.BoolVar(&helpModel, "help", false, "显示帮助")
	flag.BoolVar(&initModel, "init", false, "创建模版")
	flag.BoolVar(&indexModel, "make-index", false, "更新索引")
	flag.BoolVar(&editModel, "e", false, "编辑或新建笔记")
	flag.BoolVar(&removeModel, "r", false, "移除笔记")
	flag.BoolVar(&confirmModel, "y", false, "确认操作")
	flag.Parse()

	i = input{
		helpModel:   helpModel,
		initModel:   initModel,
		indexModel:  indexModel,
		editModel:   editModel,
		removeModel: removeModel,
		confirm:     confirmModel,
		args:        flag.Args(),
	}

	fmt.Printf("%+v\n", i)
}

func main() {

	p, err := parser(i)
	if err != nil {
		fmt.Printf("parser 有错误, 设置为help模式%v", err)
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
