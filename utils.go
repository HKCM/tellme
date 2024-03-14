package main

import (
	"os"
)

const (
	ERR      = -1
	IsNoFile = 0
	IsAFile  = 1
	IsADir   = 2
)

// 0表示指定路径不存在
// 1表示路径是文件
// 2表示路径是目录
func getPathStat(path string) (int, error) {
	e, err := pathExists(path)
	if err != nil {
		return ERR, err
	}

	if !e {
		return IsNoFile, nil
	}

	if IsFile(path) {
		return IsAFile, nil
	}

	return IsADir, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		//  no such file or directory
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
