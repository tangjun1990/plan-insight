package utilstr

import (
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
