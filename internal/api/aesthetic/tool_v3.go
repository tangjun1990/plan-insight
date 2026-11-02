package aesthetic

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cast"
	"github.com/tangjun1990/flygo/core/kcfg"
)

func getAllColorBarV3() []colorbarItem {
	return globalColorbar
}

func getAllWordV3() []string {
	allword := make([]string, 0)
	for _, item := range globalBoxV3 {
		allword = append(allword, item.Words...)
	}
	return allword
}

func getGlobalImageURLV3(img string) string {
	return kcfg.GetString("app.global.host") + "/imgv3/" + img
}

func getSolutionIconV3(boxNum int) string {
	return fmt.Sprintf("%s/solution/tag/%d.webp", kcfg.GetString("app.global.host"), boxNum)
}

func getSolutionImagesV3(boxNum int) []string {
	return []string{
		fmt.Sprintf("%s/solution/img/1/%d.webp", kcfg.GetString("app.global.host"), boxNum),
		fmt.Sprintf("%s/solution/img/2/%d.webp", kcfg.GetString("app.global.host"), boxNum),
		fmt.Sprintf("%s/solution/img/4/%d.webp", kcfg.GetString("app.global.host"), boxNum),
		fmt.Sprintf("%s/solution/img/5/%d.webp", kcfg.GetString("app.global.host"), boxNum),
		fmt.Sprintf("%s/solution/img/3/%d.webp", kcfg.GetString("app.global.host"), boxNum),
	}
}

func buildColorDescItemsV3(colorNums []int) []ColorDescItem {
	rsp := make([]ColorDescItem, 0)
	for _, num := range colorNums {
		for _, item := range globalColor {
			if item.Num != num {
				continue
			}
			words := item.Words
			if len(words) > 5 {
				words = words[:5]
			}
			rsp = append(rsp, ColorDescItem{Num: cast.ToString(item.Num), Name: item.Name, Words: words, RGB: rgbToHex(uint8(item.R), uint8(item.G), uint8(item.B))})
		}
	}
	return rsp
}

func getColorBarDetailsV3(nums []int) []ColorbarDetailV3 {
	rsp := make([]ColorbarDetailV3, 0)
	for _, num := range nums {
		for _, item := range getAllColorBarV3() {
			if item.Num == num {
				rsp = append(rsp, ColorbarDetailV3{Num: item.Num, URL: kcfg.GetString("app.global.host") + item.URL})
			}
		}
	}
	return rsp
}

func getBoxNumFromImageNameV3(imageName string) int {
	parts := strings.Split(imageName, "-")
	if len(parts) < 3 {
		return 0
	}
	return cast.ToInt(strings.Split(parts[2], ".")[0])
}

func buildAestheticSummaryV3(likedColors []int, likedColorBars []int, likedAdjectives []string, likedImages []string) (string, boxItemV3, boxItemV3, []string) {
	ensureRenderIndexesV3()
	boxScores := make(map[int]int)
	wordScores := make(map[string]int)
	for _, colorNum := range likedColors {
		for _, colorItem := range globalColor {
			if colorItem.Num == colorNum {
				for _, word := range colorItem.Words { wordScores[word]++ }
			}
		}
		if boxNum, ok := colorToBoxNumV3[colorNum]; ok { boxScores[boxNum]++ }
	}
	for _, word := range likedAdjectives {
		wordScores[word]++
		if boxNum, ok := wordToBoxNumV3[word]; ok { boxScores[boxNum]++ }
	}
	for _, imageName := range likedImages {
		if boxNum := getBoxNumFromImageNameV3(imageName); boxNum > 0 { boxScores[boxNum]++ }
	}
	for _, colorBarNum := range likedColorBars {
		if boxNum, ok := colorbarToBoxNumV3[colorBarNum]; ok { boxScores[boxNum]++ }
	}
	items := make([]boxItemV3, 0, len(boxScores))
	for boxNum := range boxScores { if box, ok := boxByNumV3[boxNum]; ok { items = append(items, box) } }
	sort.Slice(items, func(i, j int) bool { if boxScores[items[i].Num] == boxScores[items[j].Num] { return items[i].name < items[j].name }; return boxScores[items[i].Num] > boxScores[items[j].Num] })
	primary, secondary := boxItemV3{}, boxItemV3{}
	if len(items) > 0 { primary = items[0] }
	if len(items) > 1 { secondary = items[1] }
	words := make([]string, 0, len(wordScores))
	for word := range wordScores { words = append(words, word) }
	sort.Slice(words, func(i, j int) bool { if wordScores[words[i]] == wordScores[words[j]] { return words[i] < words[j] }; return wordScores[words[i]] > wordScores[words[j]] })
	if len(words) > 5 { words = words[:5] }
	if primary.Num == 0 { return "", primary, secondary, words }
	if secondary.Num == 0 { return fmt.Sprintf("“ %s派 ”", primary.name), primary, secondary, words }
	return fmt.Sprintf("“ %s%s派 ”", primary.name, secondary.name), primary, secondary, words
}
