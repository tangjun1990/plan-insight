package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cast"
)

/*type colorItem struct {
	Num     int      // 色号:1-130
	R       int      // R值
	G       int      // G值
	B       int      // B值
	Name    string   // 颜色名
	Words   []string // 对应的n个形容词
	Comment string   // 对应的文字评价
}


var globalColor = []colorItem{
	{
		1,
		192,
		6,
		32,
		"R/V大红色",
		[]string{"进取的", "精力旺盛的", "大胆的", "跃动的"},
		"红色的味感形象是辛辣、刺激、味浓有充实感。而且容易让人联想到血、火焰和太阳,因此,红色与日常生活中的吉庆、诅咒、肮脏等形象密切相连。如果同生硬的黑色或蓝色搭配,会呈现活泼,有生气的形象效果,是大胆创新的配色。",
	},
}*/

func main() {
	genWordNew()
}

/*func genword() {
	colorcontent, err := os.ReadFile("color.txt")
	if err != nil {
		log.Fatal(err)
	}
	colortxtstring := string(colorcontent)

	wordcontent, err := os.ReadFile("word.txt")
	if err != nil {
		log.Fatal(err)
	}
	wordtxtstring := string(wordcontent)
	wordSlice := strings.Split(wordtxtstring, "\n\n")
	colorWordMap := make(map[int][]string, 0)
	for k, v := range wordSlice {
		colorWordMap[k+1] = strings.Split(v, "\n")[1:]
	}

	commentcontent, err := os.ReadFile("comment.txt")
	if err != nil {
		log.Fatal(err)
	}
	commenttxtstring := string(commentcontent)
	commentSlice := strings.Split(commenttxtstring, "\n")
	validCommentSlice := make([]string, 0)
	for _, v := range commentSlice {
		if strings.Index(v, "色彩解析：") == 0 {
			validCommentSlice = append(validCommentSlice, v)
		}
	}

	colorSlice := strings.Split(colortxtstring, "\n")

	colorItemSlice := make([]colorItem, 0)
	for _, v := range colorSlice {
		tmp := strings.Split(v, ",")
		tmpNum := cast.ToInt(tmp[0])
		colorItemSlice = append(colorItemSlice, colorItem{
			Num:     tmpNum,
			Name:    tmp[1],
			R:       cast.ToInt(tmp[2]),
			G:       cast.ToInt(tmp[3]),
			B:       cast.ToInt(tmp[4]),
			Words:   colorWordMap[tmpNum],
			Comment: validCommentSlice[tmpNum-1],
		})
	}
	b3, _ := json.Marshal(colorItemSlice)
	fmt.Println(string(b3))
}*/

var WordTemplate = `
{
		%d,
		"可爱",
		0,
		0,
		[]string{%s},
		[]int{%s},
		"",
	},
`

func genWordNew() {
	wordcontent, err := os.ReadFile("word_new.txt")
	if err != nil {
		log.Fatal(err)
	}
	wordtxtstring := string(wordcontent)

	commentSlice := strings.Split(wordtxtstring, "\n")

	wordmap := make(map[int][]string, 0)
	colormap := make(map[int][]int, 0)
	for _, v := range commentSlice {

		tmp := strings.Split(v, ",")
		if len(tmp) < 3 {
			continue
		}

		curcolor := cast.ToInt(tmp[0])
		curword := tmp[1]
		curBoxNUm := cast.ToInt(tmp[2])

		if _, ok := wordmap[curBoxNUm]; !ok {
			wordmap[curBoxNUm] = make([]string, 0)
		}
		wordmap[curBoxNUm] = append(wordmap[curBoxNUm], curword)

		if _, ok := colormap[curBoxNUm]; !ok {
			colormap[curBoxNUm] = make([]int, 0)
		}
		colormap[curBoxNUm] = append(colormap[curBoxNUm], curcolor)

	}

	resultstring := ""
	for i := 1; i <= 40; i++ {
		curwordstring := ""
		for _, v := range wordmap[i] {
			curwordstring += fmt.Sprintf("\"%s\",", v)
		}
		curcolorstring := ""
		for _, v := range colormap[i] {
			curcolorstring += fmt.Sprintf("%d,", v)
		}
		resultstring += fmt.Sprintf(WordTemplate, i, strings.TrimRight(curwordstring, ","), strings.TrimRight(curcolorstring, ","))
	}
	fmt.Println(resultstring)
}
