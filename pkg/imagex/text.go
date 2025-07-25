package imagex

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// DrawTextOnImage 在图片上绘制中文文字
// imagePath: 输入图片路径
// x: 文字左上角x坐标
// y: 文字左上角y坐标（基线位置）
// fontSize: 字体大小
// text: 要绘制的文字内容
// r,g,b: 背景矩形颜色RGB值
// fontPath: 字体文件路径（.ttf文件），如果为空则使用默认字体
// outputPath: 输出图片路径，如果为空则覆盖原图片
// 返回：错误信息
func DrawTextOnImage(imagePath string, x, y int, fontSize float64, text string, r, g, b uint8, fontPath string, outputPath string) error {
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

	// 创建一个RGBA图像用于绘制
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// 加载字体
	var f *truetype.Font
	if fontPath != "" {
		// 使用指定的字体文件
		fontBytes, err := ioutil.ReadFile(fontPath)
		if err != nil {
			return err
		}
		f, err = truetype.Parse(fontBytes)
		if err != nil {
			return err
		}
	} else {
		// 使用默认字体（需要系统中有中文字体支持）
		// 这里可以尝试加载系统默认中文字体
		// 为了简化，我们使用基础字体，但可能不支持中文
		// 建议用户提供字体文件路径
		return errors.New("fontPath is required for Chinese text rendering")
	}

	// 创建字体面用于测量文字尺寸
	face := truetype.NewFace(f, &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	// 计算文字的尺寸
	textWidth := 0
	for _, char := range text {
		advance, _ := face.GlyphAdvance(char)
		textWidth += int(advance >> 6)
	}
	textHeight := int(fontSize)

	// 添加内边距
	padding := int(fontSize * 0.3) // 内边距为字体大小的30%
	rectWidth := textWidth + padding*2
	rectHeight := textHeight + padding*2

	// 计算圆角矩形的位置
	rectX := x - padding
	rectY := y - padding

	// 绘制圆角矩形背景
	cornerRadius := int(fontSize * 0.2) // 圆角半径为字体大小的20%
	drawRoundedRectangle(rgba, rectX, rectY, rectWidth, rectHeight, cornerRadius, color.RGBA{R: r, G: g, B: b, A: 255})

	// 创建freetype上下文用于绘制文字
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	// 设置文字颜色为白色
	c.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255}))
	c.SetHinting(font.HintingNone)

	// 绘制文字
	pt := freetype.Pt(x, y+int(c.PointToFixed(fontSize)>>6))
	_, err = c.DrawString(text, pt)
	if err != nil {
		return err
	}

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

// drawRoundedRectangle 绘制圆角矩形
func drawRoundedRectangle(img *image.RGBA, x, y, width, height, radius int, c color.Color) {
	// 确保圆角半径不超过矩形的一半
	if radius > width/2 {
		radius = width / 2
	}
	if radius > height/2 {
		radius = height / 2
	}

	// 绘制矩形主体（不包括圆角部分）
	// 上边矩形
	for cy := y + radius; cy < y+height-radius; cy++ {
		for cx := x; cx < x+width; cx++ {
			img.Set(cx, cy, c)
		}
	}

	// 左右边矩形
	for cy := y; cy < y+height; cy++ {
		for cx := x + radius; cx < x+width-radius; cx++ {
			img.Set(cx, cy, c)
		}
	}

	// 绘制四个圆角
	// 左上角
	drawQuarterCircle(img, x+radius, y+radius, radius, c, 2)
	// 右上角
	drawQuarterCircle(img, x+width-radius-1, y+radius, radius, c, 1)
	// 左下角
	drawQuarterCircle(img, x+radius, y+height-radius-1, radius, c, 3)
	// 右下角
	drawQuarterCircle(img, x+width-radius-1, y+height-radius-1, radius, c, 0)
}

// drawQuarterCircle 绘制四分之一圆
// quadrant: 0=右下, 1=右上, 2=左上, 3=左下
func drawQuarterCircle(img *image.RGBA, centerX, centerY, radius int, c color.Color, quadrant int) {
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			dist := dx*dx + dy*dy
			if dist <= radius*radius {
				var px, py int
				switch quadrant {
				case 0: // 右下
					if dx >= 0 && dy >= 0 {
						px, py = centerX+dx, centerY+dy
					}
				case 1: // 右上
					if dx >= 0 && dy <= 0 {
						px, py = centerX+dx, centerY+dy
					}
				case 2: // 左上
					if dx <= 0 && dy <= 0 {
						px, py = centerX+dx, centerY+dy
					}
				case 3: // 左下
					if dx <= 0 && dy >= 0 {
						px, py = centerX+dx, centerY+dy
					}
				}
				if px >= 0 && py >= 0 && px < img.Bounds().Dx() && py < img.Bounds().Dy() {
					img.Set(px, py, c)
				}
			}
		}
	}
}

// DrawMultiLineTextOnImage 在图片上绘制多行中文文字
// imagePath: 输入图片路径
// x: 文字左上角x坐标
// y: 文字左上角y坐标
// fontSize: 字体大小
// lines: 要绘制的文字行数组
// lineSpacing: 行间距
// r,g,b: 文字颜色RGB值
// fontPath: 字体文件路径（.ttf文件）
// outputPath: 输出图片路径，如果为空则覆盖原图片
// 返回：错误信息
func DrawMultiLineTextOnImage(imagePath string, x, y int, fontSize float64, lines []string, lineSpacing int, r, g, b uint8, fontPath string, outputPath string) error {
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

	// 创建一个RGBA图像用于绘制
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// 加载字体
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return err
	}

	// 创建freetype上下文
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(color.RGBA{R: r, G: g, B: b, A: 255}))
	c.SetHinting(font.HintingNone)

	// 绘制每一行文字
	for i, line := range lines {
		currentY := y + i*(int(fontSize)+lineSpacing)
		pt := freetype.Pt(x, currentY+int(c.PointToFixed(fontSize)>>6))
		_, err = c.DrawString(line, pt)
		if err != nil {
			return err
		}
	}

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
