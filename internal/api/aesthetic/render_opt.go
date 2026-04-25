package aesthetic

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/tangjun1990/plan-insight/pkg/imagex"
	"golang.org/x/image/font"
)

var (
	renderIndexOnce sync.Once
	boxByNum        map[int]boxItem
	colorToBoxNum   map[int]int
	wordToBoxNum    map[string]int
	fontCache       sync.Map
)

func ensureRenderIndexes() {
	renderIndexOnce.Do(func() {
		boxByNum = make(map[int]boxItem, len(globalBox))
		colorToBoxNum = make(map[int]int)
		wordToBoxNum = make(map[string]int)

		for _, box := range globalBox {
			boxByNum[box.Num] = box
			for _, colorNum := range box.Colors {
				colorToBoxNum[colorNum] = box.Num
			}
			for _, word := range box.Words {
				wordToBoxNum[word] = box.Num
			}
		}
	})
}

func loadImageAsRGBA(filePath string) (*image.RGBA, error) {
	img, err := GetImageFromFile(filePath)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba, nil
}

func saveRenderedImage(outputPath string, img image.Image) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	ext := strings.ToLower(filepath.Ext(outputPath))
	if ext == ".png" {
		return png.Encode(outFile, img)
	}

	return jpeg.Encode(outFile, img, &jpeg.Options{Quality: 75})
}

func drawSolidRectangleRGBA(img *image.RGBA, x, y, width, height int, c color.Color) {
	bounds := img.Bounds()
	if x < bounds.Min.X {
		x = bounds.Min.X
	}
	if y < bounds.Min.Y {
		y = bounds.Min.Y
	}
	if x+width > bounds.Max.X {
		width = bounds.Max.X - x
	}
	if y+height > bounds.Max.Y {
		height = bounds.Max.Y - y
	}
	if width <= 0 || height <= 0 {
		return
	}

	for cy := y; cy < y+height; cy++ {
		for cx := x; cx < x+width; cx++ {
			img.Set(cx, cy, c)
		}
	}

	p1 := image.Point{X: x, Y: y}
	p2 := image.Point{X: x + width - 1, Y: y}
	p3 := image.Point{X: x + width - 1, Y: y + height - 1}
	p4 := image.Point{X: x, Y: y + height - 1}
	imagex.DrawLine(img, p1, p2, c)
	imagex.DrawLine(img, p2, p3, c)
	imagex.DrawLine(img, p3, p4, c)
	imagex.DrawLine(img, p4, p1, c)
}

func drawCrossRGBA(img *image.RGBA, x, y, size, thickness int, c color.Color) {
	bounds := img.Bounds()
	if x < bounds.Min.X {
		x = bounds.Min.X
	}
	if y < bounds.Min.Y {
		y = bounds.Min.Y
	}
	if x+size > bounds.Max.X {
		size = bounds.Max.X - x
	}
	if y+size > bounds.Max.Y {
		size = bounds.Max.Y - y
	}
	if size <= 0 {
		return
	}
	if thickness <= 0 {
		thickness = 1
	}
	if thickness > size/2 {
		thickness = size / 2
	}

	for i := 0; i < thickness; i++ {
		p1 := image.Point{X: x + i, Y: y}
		p2 := image.Point{X: x + size, Y: y + size - i}
		imagex.DrawLine(img, p1, p2, c)

		if i > 0 {
			p1Alt := image.Point{X: x, Y: y + i}
			p2Alt := image.Point{X: x + size - i, Y: y + size}
			imagex.DrawLine(img, p1Alt, p2Alt, c)
		}

		p3 := image.Point{X: x + size - i, Y: y}
		p4 := image.Point{X: x, Y: y + size - i}
		imagex.DrawLine(img, p3, p4, c)

		if i > 0 {
			p3Alt := image.Point{X: x + size, Y: y + i}
			p4Alt := image.Point{X: x + i, Y: y + size}
			imagex.DrawLine(img, p3Alt, p4Alt, c)
		}
	}
}

func loadCachedFont(fontPath string) (*truetype.Font, error) {
	if cachedFont, ok := fontCache.Load(fontPath); ok {
		return cachedFont.(*truetype.Font), nil
	}

	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	parsedFont, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	actualFont, _ := fontCache.LoadOrStore(fontPath, parsedFont)
	return actualFont.(*truetype.Font), nil
}

func drawLabelTextRGBA(img *image.RGBA, x, y int, fontSize float64, text string, bgColor color.RGBA, fontPath string) error {
	parsedFont, err := loadCachedFont(fontPath)
	if err != nil {
		return err
	}

	face := truetype.NewFace(parsedFont, &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	defer face.Close()

	textWidth := 0
	for _, char := range text {
		advance, _ := face.GlyphAdvance(char)
		textWidth += int(advance >> 6)
	}
	textHeight := int(fontSize)

	padding := int(fontSize * 0.3)
	rectWidth := textWidth + padding*2
	rectHeight := textHeight + padding*2
	rectX := x - padding
	rectY := y - padding

	drawRoundedRectangleRGBA(img, rectX, rectY, rectWidth, rectHeight, int(fontSize*0.2), bgColor)

	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(parsedFont)
	ctx.SetFontSize(fontSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255}))
	ctx.SetHinting(font.HintingNone)

	pt := freetype.Pt(x, y+int(ctx.PointToFixed(fontSize)>>6))
	_, err = ctx.DrawString(text, pt)
	return err
}

func drawRoundedRectangleRGBA(img *image.RGBA, x, y, width, height, radius int, c color.Color) {
	if radius > width/2 {
		radius = width / 2
	}
	if radius > height/2 {
		radius = height / 2
	}

	for cy := y + radius; cy < y+height-radius; cy++ {
		for cx := x; cx < x+width; cx++ {
			if image.Pt(cx, cy).In(img.Bounds()) {
				img.Set(cx, cy, c)
			}
		}
	}

	for cy := y; cy < y+height; cy++ {
		for cx := x + radius; cx < x+width-radius; cx++ {
			if image.Pt(cx, cy).In(img.Bounds()) {
				img.Set(cx, cy, c)
			}
		}
	}

	drawQuarterCircleRGBA(img, x+radius, y+radius, radius, c, 2)
	drawQuarterCircleRGBA(img, x+width-radius-1, y+radius, radius, c, 1)
	drawQuarterCircleRGBA(img, x+radius, y+height-radius-1, radius, c, 3)
	drawQuarterCircleRGBA(img, x+width-radius-1, y+height-radius-1, radius, c, 0)
}

func drawQuarterCircleRGBA(img *image.RGBA, centerX, centerY, radius int, c color.Color, quadrant int) {
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy > radius*radius {
				continue
			}

			var (
				px int
				py int
				ok bool
			)

			switch quadrant {
			case 0:
				ok = dx >= 0 && dy >= 0
			case 1:
				ok = dx >= 0 && dy <= 0
			case 2:
				ok = dx <= 0 && dy <= 0
			case 3:
				ok = dx <= 0 && dy >= 0
			}

			if !ok {
				continue
			}

			px, py = centerX+dx, centerY+dy
			if image.Pt(px, py).In(img.Bounds()) {
				img.Set(px, py, c)
			}
		}
	}
}
