package golangutils

import "image"

func VerticalFlip(img image.Image) image.Image {
	// 获取原图像的宽度和高度
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建一个新的 RGBA 图像
	flipped := image.NewRGBA(bounds)

	// 遍历每个像素，并垂直翻转
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 从原图中获取像素
			originalPixel := img.At(x, y)
			// 将像素设置到翻转后的图像中
			flipped.Set(x, height-1-y, originalPixel)
		}
	}

	return flipped
}
