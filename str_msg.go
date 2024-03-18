package main

import (
	"fmt"
	"strings"
)

func deleteFileMsg(p Parser, filePath string) string {
	t := fmt.Sprintf("确认删除 %s 笔记,请使用 -y 参数: \n tellme -r", filePath)
	for _, s := range p.Args {
		t = t + " " + s
	}
	return t
}

func usage() string {
	return strings.TrimSpace(`Usage:
  tellme [option] name

Options:
  --init			create default note template
  --update-index    update-index
  -e                Edit or Create new note
  -r                Remove a note or dir
  --help            show usage

Examples:

  To show the note:
    tellme vim open
    tellme aws ec2 runinstance

  To edit (or create) the foo:
    tellme -e aws ec2 runinstance

`)
}
