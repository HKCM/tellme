package main

import (
	"fmt"
	"os"
)

type Model string

const (
	NoModel     = "NoModel"
	HelpModel   = "HelpModel"
	NormalModel = "NormalModel"
	EditModel   = "EditModel"
	RemoveModel = "RemoveModel"
)

type Parser struct {
	model      Model
	parameters []string
}

func parser() (p Parser, err error) {

	p.model = NormalModel
	// 直接运行命令没有参数 直接返回help
	if len(os.Args) == 1 {
		fmt.Println("help model")
		p.model = HelpModel
		return
	}

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h", "help": // 如果是help模式直接返回
			p.model = HelpModel
			return
		case "--edit", "-e":
			err = setModel(&p, EditModel)
		case "--remove", "-r":
			err = setModel(&p, RemoveModel)
		default:
			p.parameters = append(p.parameters, arg)
		}
		if err != nil {
			return p, err
		}
	}

	lenP := len(p.parameters)
	// 参数最多为三个 例如 tellme aws ec2 runinstance
	if lenP > Deep {
		p.model = HelpModel
	}

	// 如果不是帮助模式,但是关键词为0 则显示帮助
	if p.model != HelpModel && lenP == 0 {
		p.model = HelpModel
	}

	return
}

func setModel(p *Parser, model Model) error {
	if p.model == NormalModel {
		p.model = model
		return nil
	} else {
		return fmt.Errorf("不能同时指定多个模式: %s,%s", p.model, model)
	}
}
