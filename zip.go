package golangutils

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// 如果是目录，则创建目录
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// 创建文件的目录
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// 创建文件
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// 关闭文件和读取器
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func IsCompressedFile(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	// 读取文件的前4个字节
	buffer := make([]byte, 4)
	if _, err := file.Read(buffer); err != nil {
		return false
	}

	// 检查魔数
	switch {
	case bytes.HasPrefix(buffer, []byte{0x50, 0x4B, 0x03, 0x04}): // ZIP
		return true
	case bytes.HasPrefix(buffer, []byte{0x1F, 0x8B}): // GZIP
		return true
	// 添加其他格式的魔数检查...
	default:
		return false
	}
}
