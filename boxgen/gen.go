package main

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"

	"git.4321.sh/feige/commonapi/internal/api/aesthetic"
	"git.4321.sh/feige/commonapi/pkg/imagex"
	"github.com/disintegration/imaging"
	"github.com/spf13/cast"
)

func main() {
	drawImageToBox()
}

// 将指定的多张图片按对应的坐标startX和startY放进图片boxbase.jpg中
func drawImageToBox() {
	imgs := []string{"1-1-1.png", "2-1-1.jpg", "3-1-1.png", "4-1-1.png", "1-1-2.png", "1-1-3.png", "1-1-4.png", "1-1-5.png", "1-1-6.png", "1-1-7.png", "1-1-8.png", "1-1-9.png", "1-1-10.png", "1-1-11.png", "1-1-12.png", "1-1-13.png", "1-1-14.png", "1-1-15.png", "1-1-16.png"}
	//imgs := []string{"1-1-9.png"}
	backImagePath := "./boxbase-3.jpg"
	backgroundImage, _ := GetImageFromFile(backImagePath)

	// 输出图片路径
	outputPath := "./out_example_2.jpg"

	//words := []string{"可爱的", "雅致的", "快乐的", "进取的", "坚韧的", "甜美的"}
	likedColor := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	dislikedColor := []int{10, 20, 30, 40, 50}

	// 1.把图片合并到抽屉图
	var overlap *image.NRGBA
	imagemap := make(map[int][]string, 0)
	for _, img := range imgs {
		tmp := strings.Split(img, "-")
		tmpSuffix := tmp[2]
		tmpNum := strings.Split(tmpSuffix, ".")
		boxNum := cast.ToInt(tmpNum[0])

		startX := 0
		startY := 0
		for _, vv := range globalBox {
			if vv.Num == boxNum {
				startX = vv.StartX
				startY = vv.StartY
			}
		}
		// 如果这个box内已有图片了，startx需要向右移动200
		if _, ok := imagemap[boxNum]; ok {
			startX = startX + (len(imagemap[boxNum]) * 200)
		}

		srcImage, _ := GetImageFromFile(img)
		overlap = imaging.Overlay(backgroundImage, srcImage, image.Point{startX, startY}, 0.75)
		backgroundImage = overlap
		if _, ok := imagemap[boxNum]; ok {
			imagemap[boxNum] = append(imagemap[boxNum], img)
		} else {
			imagemap[boxNum] = []string{img}
		}
	}
	_ = imaging.Save(overlap, outputPath)
	// 2.把关键词合并到抽屉图
	/*wordmap := make(map[int][]string, 0)
	for _, v := range words {
		for k, vv := range globalBox {
			for _, vvv := range vv.Words {
				if vvv == v {
					if _, ok := wordmap[vv.Num]; ok {
						wordmap[vv.Num] = append(wordmap[vv.Num], v)
					} else {
						wordmap[vv.Num] = []string{v}
					}
				}
			}
		}
	}
	for _, v := range wordmap {
		startX := 0
		startY := 0
		for _, vv := range globalBox {
			if vv.Num == boxNum {
				startX = vv.StartX
				startY = vv.StartY
			}
		}
		tmpstring := strings.Join(v, ",")
		// 使用freetype将文字写入图片overlap中
	}
	*/

	// 3.把喜欢的颜色合并到抽屉图
	likedcolormap := make(map[int][]int, 0)
	for _, v := range likedColor {
		for _, vv := range globalBox {
			for _, vvv := range vv.Colors {
				if vvv == v {
					if _, ok := likedcolormap[vv.Num]; ok {
						likedcolormap[vv.Num] = append(likedcolormap[vv.Num], v)
					} else {
						likedcolormap[vv.Num] = []int{v}
					}
				}
			}
		}
	}
	for _, v := range likedcolormap {
		startX := 0
		startY := 0
		for _, col := range v {
			for _, vv := range globalBox {
				for _, vvv := range vv.Colors {
					if vvv == col {
						startX = vv.StartX
						startY = vv.StartY
					}
				}
			}
		}
		for _, col := range v {
			// 画喜欢的颜色，x不变，y增加250
			curY := startY + 250
			r, g, b := aesthetic.NumToRGB(col)
			// 画喜欢的颜色，使用矩形
			err := imagex.DrawRectangleOnImage(outputPath, startX, curY, 95, 93, r, g, b, outputPath)
			if err != nil {

			}
			// 如果是多个颜色，x增加100
			startX = startX + 100
		}
	}

	// 4.把不喜欢的颜色合并到抽屉图
	dislikedcolormap := make(map[int][]int, 0)
	for _, v := range dislikedColor {
		for _, vv := range globalBox {
			for _, vvv := range vv.Colors {
				if vvv == v {
					if _, ok := dislikedcolormap[vv.Num]; ok {
						dislikedcolormap[vv.Num] = append(dislikedcolormap[vv.Num], v)
					} else {
						dislikedcolormap[vv.Num] = []int{v}
					}
				}
			}
		}
	}
	for _, v := range dislikedcolormap {
		startX := 0
		startY := 0
		for _, col := range v {
			for _, vv := range globalBox {
				for _, vvv := range vv.Colors {
					if vvv == col {
						startX = vv.StartX
						startY = vv.StartY
					}
				}
			}
		}
		for _, col := range v {
			// 画不喜欢的颜色，x不变，y增加250
			curY := startY + 350
			r, g, b := aesthetic.NumToRGB(col)
			// 画不喜欢的颜色，使用矩形
			err := imagex.DrawRectangleOnImage(outputPath, startX, curY, 50, 93, r, g, b, outputPath)
			if err != nil {

			}
			// 如果是多个颜色，x增加55
			startX = startX + 55
		}
	}

	// 5. 保存图片
	//_ = imaging.Save(overlap, outputPath)
}

func GetImageFromFile(filePath string) (img image.Image, err error) {
	f1Src, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer f1Src.Close()

	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = f1Src.Read(buff)

	if err != nil {
		return nil, err
	}

	filetype := http.DetectContentType(buff)

	fSrc, err := os.Open(filePath)
	defer fSrc.Close()

	switch filetype {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(fSrc)
		if err != nil {
			return nil, err
		}

	case "image/gif":
		img, err = gif.Decode(fSrc)
		if err != nil {
			return nil, err
		}

	case "image/png":
		img, err = png.Decode(fSrc)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	return img, nil
}

type boxItem struct {
	Num    int
	name   string
	StartX int
	StartY int
	Words  []string
	Colors []int
}

// 以下数据适用于5000*3500的base图片
var globalBox = []boxItem{
	{
		1,
		"可爱的",
		1210,
		510,
		[]string{"稚嫩的", "可爱的", "孩子气的", "伶俐的"},
		[]int{2, 3, 21, 22, 30, 31, 23},
	},
	{
		2,
		"闲适的",
		1050,
		1140,
		[]string{"开朗的", "快乐的", "高兴的", "愉快的", "风趣的", "阳光的", "快活的", "活跃的", "有生气的", "闲适的", "朝气蓬勃的", "鲜艳的", "绚丽的", "轻松的", "自由自在的", "无忧无虑的", "有亲和力的", "开放的"},
		[]int{1, 10, 11, 12},
	},
	{
		3,
		"动感的",
		1050,
		2050,
		[]string{"充满活力的", "悦动的", "进取的", "主动的", "大胆的", "刺激的", "热烈的", "激烈的", "强烈的", "动感的", "力动的", "精力旺盛的"},
		[]int{0},
	},
	{
		4,
		"豪华的",
		1420,
		2050,
		[]string{"娇媚的", "娇艳的", "华丽的", "性感的", "魅惑的", "富于装饰的", "丰满的", "丰润的", "豪华的", "成熟的", "奢华的", "充实的", "浓郁的"},
		[]int{91, 92, 100},
	},
	{
		5,
		"粗犷的",
		1420,
		2730,
		[]string{"强劲的", "坚韧的", "阳刚的", "健壮的"},
		[]int{101},
	},

	{
		6,
		"浪漫的",
		1840,
		430,
		[]string{"楚楚动人的", "甜美的", "纯净的", "浪漫的", "童话般的", "朦胧的", "纯真的", "清纯的"},
		[]int{32, 33, 34, 41, 42, 43, 44, 49, 50, 122},
	},
	{
		7,
		"自然的",
		1840,
		820,
		[]string{"自然的", "悠然自然的", "温和的", "大方的", "放松的", "舒适的", "坦诚的", "悠闲的", "家居的", "温润的", "和睦的", "柔软的", "柔和的",
			"融洽的", "温柔的", "安宁的", "温顺的", "淡泊的", "简朴的", "不加修饰的", "和平的", "水灵灵的", "惬意的", "健康的", "新鲜的", "鲜活的", "安稳的"},
		[]int{13, 20, 24, 25, 39, 40, 51, 52, 53, 54, 55, 59, 61, 62, 70, 82, 123},
	},
	{
		8,
		"雅致的",
		2310,
		1060,
		[]string{"优美的", "有情趣的", "端庄的", "抒情的", "细腻的", "细致的", "柔美的", "娇美的", "有品味位的", "含蓄的", "华美的", "女性化的", "雅致的", "温文尔雅的", "秀丽的", "优雅的"},
		[]int{4, 5, 14, 19, 29, 60, 63, 64, 69, 83, 84, 124, 125},
	},
	{
		9,
		"精致的",
		2930,
		1250,
		[]string{"微妙的", "谨慎的", "安静的", "随章的", "洗练的", "质朴的", "精致的", "洒脱的", "江河的", "都市气息的", "文化气息的", "知性的", "冷静的", "娴静的", "萧瑟的", "素雅的", "风流的", "乡土气息的"},
		[]int{6, 13, 15, 28, 65, 66, 68, 85, 86, 87, 89},
	},
	{
		10,
		"古典的",
		2070,
		2050,
		[]string{"怀念的", "古风的", "深邃的", "潜心的", "怀旧的", "传统的", "古典的"},
		[]int{9, 71, 72, 81, 90, 93, 99, 102, 103},
	},
	{
		11,
		"考究的",
		2600,
		2210,
		[]string{"深沉的", "绅士的", "男子汉的", "严谨的", "凌然的", "正统的", "坚实的", "考究的"},
		[]int{8, 73, 74, 75, 76, 77, 78, 79, 80, 88, 94, 95, 96, 97, 98, 126, 127},
	},
	{
		12,
		"古典的&考究的",
		2160,
		2730,
		[]string{"锤炼的", "庄重的", "厚重的", "坚定的", "有格调的", "独到的", "正宗的"},
		[]int{104, 105, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 128, 129, 130},
	},
	{
		13,
		"正式的",
		3060,
		2690,
		[]string{"高雅的", "高沿的", "高贵的", "神圣的", "庄严的", "肃穆的"},
		[]int{106, 107},
	},
	{
		14,
		"清爽的",
		3290,
		610,
		[]string{"清雅的", "清新的", "清静的", "清澈的", "清爽的", "清朗的", "清冷的"},
		[]int{26, 35, 36, 37, 38, 45, 46, 47, 48, 56, 57, 58, 121},
	},
	{
		15,
		"冷.闲适的",
		3390,
		1240,
		[]string{"轻快的", "青春的", "青春洋溢的", "清冽的", "运动的"},
		[]int{16, 17, 27, 67},
	},
	{
		16,
		"现代的",
		3440,
		2070,
		[]string{"迅捷的", "现代的", "进步的", "革新的", "机敏的", "理性的", "敏锐的", "精确的", "合理的", "致密的", "人工的", "现代化的"},
		[]int{7},
	},
}

/*
以下数据适用于3000*2100的base图片
var globalBox = []boxItem{
	{
		1,
		"可爱的",
		730,
		310,
		[]string{"稚嫩的", "可爱的", "孩子气的", "伶俐的"},
	},
	{
		2,
		"闲适的",
		620,
		680,
		[]string{"开朗的", "快乐的", "高兴的", "愉快的", "风趣的", "阳光的", "快活的", "活跃的", "有生气的", "闲适的", "朝气蓬勃的", "鲜艳的", "绚丽的",
			"轻松的", "自由自在的", "无忧无虑的", "有亲和力的", "开放的"},
	},
	{
		3,
		"动感的",
		620,
		1220,
		[]string{"充满活力的", "悦动的", "进取的", "主动的", "大胆的", "刺激的", "热烈的", "激烈的", "强烈的", "动感的", "力动的", "精力旺盛的"},
	},
	{
		4,
		"豪华的",
		850,
		1220,
		[]string{"娇媚的", "娇艳的", "华丽的", "性感的", "魅惑的", "富于装饰的", "丰满的", "丰润的", "豪华的", "成熟的", "奢华的", "充实的", "浓郁的"},
	},
	{
		5,
		"粗犷的",
		850,
		1630,
		[]string{"强劲的", "坚韧的", "阳刚的", "健壮的"},
	},
	{
		6,
		"浪漫的",
		1100,
		255,
		[]string{"楚楚动人的", "甜美的", "纯净的", "浪漫的", "童话般的", "朦胧的", "纯真的", "清纯的"},
	},
	{
		7,
		"自然的",
		1100,
		490,
		[]string{"自然的", "悠然自然的", "温和的", "大方的", "放松的", "舒适的", "坦诚的", "悠闲的", "家居的", "温润的", "和睦的", "柔软的", "柔和的",
			"融洽的", "温柔的", "安宁的", "温顺的", "淡泊的", "简朴的", "不加修饰的", "和平的", "水灵灵的", "惬意的", "健康的", "新鲜的", "鲜活的", "安稳的"},
	},
	{
		8,
		"雅致的",
		1390,
		640,
		[]string{"优美的", "有情趣的", "端庄的", "抒情的", "细腻的", "细致的", "柔美的", "娇美的", "有品味位的", "含蓄的", "华美的", "女性化的", "雅致的", "温文尔雅的", "秀丽的", "优雅的"},
	},
	{
		9,
		"精致的",
		1755,
		745,
		[]string{"微妙的", "谨慎的", "安静的", "随章的", "洗练的", "质朴的", "精致的", "洒脱的", "江河的", "都市气息的", "文化气息的", "知性的", "冷静的", "娴静的", "萧瑟的", "素雅的", "风流的", "乡土气息的"},
	},
	{
		10,
		"古典的",
		1250,
		1220,
		[]string{"怀念的", "古风的", "深邃的", "潜心的", "怀旧的", "传统的", "古典的"},
	},
	{
		11,
		"考究的",
		1560,
		1320,
		[]string{"深沉的", "绅士的", "男子汉的", "严谨的", "凌然的", "正统的", "坚实的", "考究的"},
	},
	{
		12,
		"古典的&考究的",
		1300,
		1630,
		[]string{"锤炼的", "庄重的", "厚重的", "坚定的", "有格调的", "独到的", "正宗的"},
	},
	{
		13,
		"正式的",
		1830,
		1610,
		[]string{"高雅的", "高沿的", "高贵的", "神圣的", "庄严的", "肃穆的"},
	},
	{
		14,
		"清爽的",
		1980,
		355,
		[]string{"清雅的", "清新的", "清静的", "清澈的", "清爽的", "清朗的", "清冷的"},
	},
	{
		15,
		"冷.闲适的",
		2030,
		745,
		[]string{"轻快的", "青春的", "青春洋溢的", "清冽的", "运动的"},
	},
	{
		16,
		"现代的",
		2070,
		1245,
		[]string{"迅捷的", "现代的", "进步的", "革新的", "机敏的", "理性的", "敏锐的", "精确的", "合理的", "致密的", "人工的", "现代化的"},
	},
}
*/
