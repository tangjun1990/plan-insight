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

// DrawTriangleOnImage 在图片中间添加红色三角形
// imagePath: 输入图片路径
// outputPath: 输出图片路径，如果为空则覆盖原图片
// 返回：错误信息
func DrawTriangleOnImage(imagePath string, outputPath string) error {
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
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建一个RGBA图像用于绘制
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// 定义红色
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// 计算三角形顶点坐标
	// 三角形大小为图片较短边的1/4
	size := width
	if height < width {
		size = height
	}
	size = size / 4

	// 三角形顶点坐标
	centerX := width / 2
	centerY := height / 2

	p1 := image.Point{X: centerX, Y: centerY - size/2}
	p2 := image.Point{X: centerX - size/2, Y: centerY + size/2}
	p3 := image.Point{X: centerX + size/2, Y: centerY + size/2}

	// 绘制三角形 - 完全填充红色
	drawSolidTriangle(rgba, p1, p2, p3, red)

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

// 绘制实心三角形，确保完全填充
func drawSolidTriangle(img *image.RGBA, p1, p2, p3 image.Point, c color.Color) {
	// 确定三角形的边界框
	minX := Min(p1.X, Min(p2.X, p3.X))
	maxX := Max(p1.X, Max(p2.X, p3.X))
	minY := Min(p1.Y, Min(p2.Y, p3.Y))
	maxY := Max(p1.Y, Max(p2.Y, p3.Y))

	// 对边界框内的每个像素，判断是否在三角形内部
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 判断点(x,y)是否在三角形内部
			if pointInTriangle(image.Point{X: x, Y: y}, p1, p2, p3) {
				img.Set(x, y, c)
			}
		}
	}

	// 绘制三角形边缘以确保边缘像素被填充
	DrawLine(img, p1, p2, c)
	DrawLine(img, p2, p3, c)
	DrawLine(img, p3, p1, c)
}

// 判断点p是否在由p1,p2,p3确定的三角形内部
func pointInTriangle(p, p1, p2, p3 image.Point) bool {
	// 使用叉积判断点是否在三角形内部
	// 如果点p在三角形p1,p2,p3内部，则p在p1p2, p2p3, p3p1三条边的同一侧

	// 计算叉积 (p1-p3)×(p-p3)
	v1x, v1y := p1.X-p3.X, p1.Y-p3.Y
	v2x, v2y := p.X-p3.X, p.Y-p3.Y
	c1 := v1x*v2y - v1y*v2x

	// 计算叉积 (p2-p1)×(p-p1)
	v1x, v1y = p2.X-p1.X, p2.Y-p1.Y
	v2x, v2y = p.X-p1.X, p.Y-p1.Y
	c2 := v1x*v2y - v1y*v2x

	// 如果c1和c2同号（都为正或都为负）
	if (c1 > 0 && c2 < 0) || (c1 < 0 && c2 > 0) {
		return false
	}

	// 计算叉积 (p3-p2)×(p-p2)
	v1x, v1y = p3.X-p2.X, p3.Y-p2.Y
	v2x, v2y = p.X-p2.X, p.Y-p2.Y
	c3 := v1x*v2y - v1y*v2x

	// 判断c3与c1、c2同号
	return !((c3 > 0 && c1 < 0) || (c3 < 0 && c1 > 0))
}
