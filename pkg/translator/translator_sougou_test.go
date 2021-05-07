package translator

import (
	"testing"
)

func TestSougouTranslate(t *testing.T) {

	tl, ok := ENGINES[EngineSougou]
	cases := []struct {
		from, to, text string
	}{
		{"en", "zh", "Did you have eaten?"},
		{"zh", "en", "你吃饭了没有?"},
		{"zh", "en", "我不知道你竟然不知道我不知道这件事。"},
	}
	if ok {
		for _, c := range cases {
			r, err := tl.Translate(c.from, c.to, c.text)
			if err == nil {
				t.Log(r.String())
			}
		}
	}
}
