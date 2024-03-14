package main

import (
	"log"
	"os/user"
	"strings"
)

const (
	editor = "vim"
)

var home string

const Deep = 3

func init() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	home = u.HomeDir + "/.tellme"
}

func main() {

	p, err := parser()
	if err != nil {
		panic(err)
	}

	// 组装路径
	var path, keyword string
	if len(p.parameters) != 0 {
		path = strings.Join(p.parameters[:len(p.parameters)-1], "/")
		keyword = p.parameters[len(p.parameters)-1]
	}

	cmd, err := selectCmd(p.model)
	if err != nil {
		panic(err)
	}

	var runCmd func(path, keyword string)

	switch cmd {
	case HelpCmd:
		runCmd = cmdShowHelp
	case EditCmd:
		runCmd = cmdEditNote
	case ShowCmd:
		runCmd = cmdShow
	case RmCmd:
		runCmd = cmdRemoveNote
	}

	runCmd(path, keyword)

}
