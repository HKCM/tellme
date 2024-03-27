package main

import (
	"strings"
)

func usage() string {
	return strings.TrimSpace(`Usage:
  tellme [option] <note|tag>
Options:
  -t                show folder tag
  -e                edit or create new note
  -r                remove a note or dir
  -y                confirm the operation
  --init            create default note template
  --update-tag      update-tag, default update all tag
  --target          use with update-tag flag
  --debug           show debug message
  --help            show usage

Examples:
  To show the note:
    tellme ln                  # 默认显示shell文件夹中笔记  
    tellme aws ec2 runinstance # 显示文件

  To edit (or create) the note:
    tellme -e aws ec2 runinstance  # 创建或编辑新笔记

  Other function
    tellme --init                      # 初始化Template文件
    tellme aws                         # 获取文件夹内文件
    tellme -t aws                      # 获取文件夹内的所有文件的tag
    tellme --update-tag                # 更新所有文件夹目录
    tellme --update-tag --target shell # 更新指定文件夹目录
`)
}
