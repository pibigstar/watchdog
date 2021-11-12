package watchdog

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func removeFileByPrefix(basePath string, prefix string) error {
	_ = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), prefix) {
			return os.Remove(filepath.Join(basePath, info.Name()))
		}
		return nil
	})
	return nil
}

// 判断文件夹是否存在，不存在则创建
// 判断文件夹是否有可写权限
func checkPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	if err := syscall.Access(path, syscall.O_RDWR); err != nil {
		return err
	}
	return nil
}
