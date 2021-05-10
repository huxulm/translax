package translator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYoudaoTranslatePrepareForm(t *testing.T) {
	cases := []struct {
		expect, text string
	}{
		{
			expect: "i=%E4%BD%A0%E7%8E%B0%E5%9C%A8%E5%93%AA%E9%87%8C%EF%BC%9F&from=AUTO&to=AUTO&smartresult=dict&client=fanyideskweb&salt=16206301042577&sign=b86db58418044ae289b2e8b7f00cdbd3&lts=1620630104257&bv=940ae85ba32dcefb13c38faaf66f115f&doctype=json&version=2.1&keyfrom=fanyi.web&action=FY_BY_REALTlME",
			text:   "你现在哪里？",
		},
	}

	for _, c := range cases {
		y := NewYoudao(nil)
		v := (y.(*youdao)).prepareForm("", "", "你现在哪里？")
		assert.Equal(t, c.expect, v.Encode())
	}
}

func TestYoudaoTranslate(t *testing.T) {

	tl, ok := ENGINES[EngineYoudao]
	cases := []struct {
		from, to, text string
	}{
		{"en", "zh", "Did you have eaten?"},
		{"zh", "en", "你吃饭了没有?"},
		{"zh", "en", "我不知道你竟然不知道我不知道这件事。"},
		{"en", "zh", longText},
	}
	if ok {
		for i := range cases {
			c := cases[i]
			r, err := tl.Translate(c.from, c.to, c.text)
			if err == nil {
				t.Log(r.String())
				// write to file
				// ioutil.WriteFile(fmt.Sprintf("./yd_%d.out", i+1), []byte(r.String()), os.ModePerm)
			}
		}
	}
}
