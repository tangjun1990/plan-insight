package imagex

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// DrawRectangleOnImage 在图片上绘制矩形
// imagePath: 输入图片路径
// x: 矩形左上角x坐标
// y: 矩形左上角y坐标
// width: 矩形宽度
// height: 矩形高度
// r,g,b: 矩形颜色RGB值
// outputPath: 输出图片路径，如果为空则覆盖原图片
// 返回：错误信息
func DrawRectangleOnImage(imagePath string, x, y, width, height int, r, g, b uint8, outputPath string) error {
	// 如果输出路径为空，则覆盖原图片
	if outputPath == "" {
		outputPath = imagePath
	}

	// 打开图片文件
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 解码图片
	var img image.Image
	ext := strings.ToLower(filepath.Ext(imagePath))

	if ext == ".png" {
		img, err = png.Decode(file)
	} else {
		// 默认按JPEG处理
		img, err = jpeg.Decode(file)
	}
	if err != nil {
		return err
	}

	// 获取图片尺寸
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// 创建一个RGBA图像用于绘制
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// 定义颜色
	rectColor := color.RGBA{R: r, G: g, B: b, A: 255}

	// 确保矩形不超出图片边界
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x+width > imgWidth {
		width = imgWidth - x
	}
	if y+height > imgHeight {
		height = imgHeight - y
	}

	// 绘制矩形
	drawSolidRectangle(rgba, x, y, width, height, rectColor)

	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 保存图片
	ext = strings.ToLower(filepath.Ext(outputPath))
	if ext == ".png" {
		err = png.Encode(outFile, rgba)
	} else {
		// 默认保存为JPEG
		err = jpeg.Encode(outFile, rgba, nil)
	}

	return err
}

// 绘制实心矩形
func drawSolidRectangle(img *image.RGBA, x, y, width, height int, c color.Color) {
	// 填充矩形内部
	for cy := y; cy < y+height; cy++ {
		for cx := x; cx < x+width; cx++ {
			img.Set(cx, cy, c)
		}
	}

	// 绘制矩形边框（虽然已经填充，但为了确保边界像素被正确填充）
	p1 := image.Point{X: x, Y: y}
	p2 := image.Point{X: x + width - 1, Y: y}
	p3 := image.Point{X: x + width - 1, Y: y + height - 1}
	p4 := image.Point{X: x, Y: y + height - 1}

	DrawLine(img, p1, p2, c)
	DrawLine(img, p2, p3, c)
	DrawLine(img, p3, p4, c)
	DrawLine(img, p4, p1, c)
}
