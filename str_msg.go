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

`)
}
