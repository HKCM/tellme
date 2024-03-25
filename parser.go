package main

import (
	"flag"
	"fmt"
	"log/slog"
	"strings"
)

type Model int

const (
	ModelHelp Model = iota
	ModelInit
	ModelUpdateKeyword
	ModelEdit
	ModelRemove
	ModelShow
)

type Parser struct {
	Model        Model
	Tag          string
	Path         string
	Confirm      bool
	ShowTags     bool
	UpdateTarget string
	Args         []string
	ArgL         int
}

type input struct {
	helpModel        bool
	showKeywordModel bool
	initModel        bool
	keywordModel     bool
	editModel        bool
	removeModel      bool
	confirm          bool
	tagTarget        string
	args             []string
}

func parser(i input) (p Parser, err error) {

	// 定义命令行参数

	// 如果没有flag并且arg为0则为帮助模式 例如直接运行命令
	// 如果arg大于指定深度则为帮助模式 例如 tellme aws ec2 runinstance
	if flag.NArg() > Deep || (flag.NFlag() == 0 && flag.NArg() == 0) {
		err = fmt.Errorf("没有flag和参数,或参数个数大于%d", Deep)
		fmt.Println(err)
		p.Model = ModelHelp
		return
	}

	models := []bool{
		i.helpModel,
		i.initModel,
		i.keywordModel,
		i.editModel,
		i.removeModel,
	}

	n := 0
	for _, b := range models {
		if b {
			n++
		}
	}

	if n > 1 {
		err = fmt.Errorf("指定多个Flag")
		fmt.Println(err)
		return
	}
	p.Confirm = i.confirm

	switch {
	case i.helpModel:
		p.Model = ModelHelp
	case i.initModel:
		p.Model = ModelInit
	case i.keywordModel:
		p.Model = ModelUpdateKeyword
	case i.editModel:
		p.Model = ModelEdit
	case i.removeModel:
		p.Model = ModelRemove
	default:
		p.Model = ModelShow
		p.ShowTags = i.showKeywordModel
	}

	p.UpdateTarget = i.tagTarget
	p.Args = i.args

	p.ArgL = len(i.args)

	switch p.ArgL {
	case 0:
		slog.Debug("No args")
	case 1:
		p.Path = DefaultPath
		p.Tag = i.args[0]
	default:
		p.Path = strings.Join(i.args[:p.ArgL-1], "/")
		p.Tag = i.args[p.ArgL-1]
	}

	return
}
