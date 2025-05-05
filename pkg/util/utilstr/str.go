package utilstr

import (
	"github.com/gofrs/uuid"
	"regexp"
	"strings"
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
