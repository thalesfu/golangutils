package golangutils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"log"
)

func ConvertToPngBase64String(img image.Image) string {
	// 将图像编码为 PNG 格式
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatalf("Failed to encode PNG: %v", err)
	}

	// 转为 Base64
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
