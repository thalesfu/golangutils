package golangutils

import (
	"fmt"
	"github.com/ftrvxmtrx/tga"
	"image"
	"os"
)

func LoadTgaFile(path string) (image.Image, bool) {
	if _, err := os.Stat(path); err != nil {
		return nil, false
	}

	bmpFile, err := os.Open(path)
	if err != nil {
		fmt.Println("无法打开TGA文件:", err)
		return nil, false
	}
	defer func(bmpFile *os.File) {
		err := bmpFile.Close()
		if err != nil {
			fmt.Println("关闭TGA文件错误:", err)
		}
	}(bmpFile)

	img, err := tga.Decode(bmpFile)
	if err != nil {
		fmt.Println("TGA解码错误:", err)
		return nil, false
	}

	return img, true
}
