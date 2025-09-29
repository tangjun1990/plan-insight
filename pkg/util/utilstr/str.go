package utilstr

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"github.com/gofrs/uuid"
)

var clearHtmlReg = regexp.MustCompile(`<[\S\s]+?>`)

// ClearHtml 清除字符串类似 "<.*>"的标签,并 TrimSpace
func ClearHtml(src string) string {
	return strings.TrimSpace(clearHtmlReg.ReplaceAllString(src, ""))
}

func UUID() string {
	str, _ := uuid.NewV4()
	return str.String()
}

// SortStrings 对字符串切片进行排序并返回排序后的结果
func SortStrings(strs []string) []string {
	// 创建一个新的切片副本，避免修改原始切片
	sorted := make([]string, len(strs))
	copy(sorted, strs)

	// 对副本进行排序
	sort.Strings(sorted)

	return sorted
}

func SortLikedImage(likedImageStr string) string {
	likedImage := make([]string, 0)
	_ = json.Unmarshal([]byte(likedImageStr), &likedImage)
	// likedImage 中元素的格式为"4-1-23.webp","1-1-13.webp","2-1-38.webp","3-1-2.webp","5-1-28.webp","6-1-25.webp"
	// 需要改为按首个数字进行排序
	sort.Slice(likedImage, func(i, j int) bool {
		// 提取第一个数字
		numI := strings.Split(likedImage[i], "-")[0]
		numJ := strings.Split(likedImage[j], "-")[0]
		return numI < numJ
	})

	return strings.Join(likedImage, ",")
}

func SortLikedImageSlice(likedImageStr []string) []string {
	// 创建一个新的切片副本，避免修改原始切片
	sorted := make([]string, len(likedImageStr))
	copy(sorted, likedImageStr)

	// 对副本进行排序
	sort.Slice(sorted, func(i, j int) bool {
		// 提取第一个数字
		numI := strings.Split(sorted[i], "-")[0]
		numJ := strings.Split(sorted[j], "-")[0]
		return numI < numJ
	})

	return sorted
}
