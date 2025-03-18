package golangutils

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"os"
)

func LoadBmpFile(path string) (image.Image, bool) {
	if _, err := os.Stat(path); err != nil {
		return nil, false
	}

	bmpFile, err := os.Open(path)
	if err != nil {
		fmt.Println("无法打开BMP文件:", err)
		return nil, false
	}
	defer func(bmpFile *os.File) {
		err := bmpFile.Close()
		if err != nil {
			fmt.Println("关闭BMP文件错误:", err)
		}
	}(bmpFile)

	img, err := bmp.Decode(bmpFile)
	if err != nil {
		fmt.Println("BMP解码错误:", err)
		return nil, false
	}

	return img, true
}
