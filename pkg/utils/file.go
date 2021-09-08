// Package utils  文件操作类
package utils

import (
	"os"
	"path/filepath"
)

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return !os.IsNotExist(err)
}

func BinBaseDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
