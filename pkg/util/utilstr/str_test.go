package utilstr

import "testing"

func TestClearHtml(t *testing.T) {
	cases := [][2]string{
		{"<p>你们医院在什么地方？<br></p>", "你们医院在什么地方？"},
		{"<p>你好<br></p>", "你好"},
		{"<p>你好  <br> </p>", "你好"},
	}

	for _, s := range cases {
		if s[1] != ClearHtml(s[0]) {
			t.Errorf("ClearHtml(%s) 期望:%s, 结果:%s", s[0], s[1], ClearHtml(s[0]))
		}
	}
}
