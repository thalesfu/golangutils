package golangutils

import (
	"golang.org/x/text/encoding"
	"log"
	"os"
	"path/filepath"
)

func LoadContentWithEncoding(path string, encode encoding.Encoding) (string, bool) {
	if _, err := os.Stat(path); err != nil {
		return "", false
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("读取文件失败：%s\nerror:%v", path, err)
		return "", false
	}

	decoder := encode.NewDecoder()

	bytes, err = decoder.Bytes(bytes)
	if err != nil {
		log.Printf("通过%T编码：%s报错\nerror:%v", encode, path, err)
		return "", false
	}

	return string(bytes), true
}

func LoadContent(path string) (string, bool) {
	if _, err := os.Stat(path); err != nil {
		return "", false
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("读取文件失败：%s\nerror:%v", path, err)
		return "", false
	}

	return string(bytes), true
}

func WriteContent(filePath string, content string) bool {
	dirPath := filepath.Dir(filePath)

	// 创建目录，包括任何必需的父目录，权限设置为 0755
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Printf("创建目录失败：%s\nerror:%v", dirPath, err)
		return false
	}

	// 创建新的文件
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("创建文件失败：%s\nerror:%v", filePath, err)
		return false
	}

	defer func(f *os.File) {
		log.Printf("[文件保存]写入并保存文件: %s", f.Name())
		err := f.Close()
		if err != nil {
			log.Printf("关闭文件失败：%s\nerror:%v", f.Name(), err)
		}
	}(file)

	// 将内容写入文件
	if _, err := file.WriteString(content); err != nil {
		log.Printf("写入文件失败：%s\nerror:%v", filePath, err)
		return false
	}

	return true
}
