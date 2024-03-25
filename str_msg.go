package main

import (
	"strings"
)

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

    go run . --debug --update-tag
    go run . --debug --update-tag --target mysql
    go run . aws # 获取文件夹内文件
    go run . -t aws # 获取文件夹内Tags
    go run . --debug --update-tag # 更新所有文件夹目录
    go run . --debug --update-tag --target shell # 更新指定文件夹目录

`)
}
