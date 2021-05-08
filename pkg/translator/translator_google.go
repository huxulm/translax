package translator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

const (
	GOOGLE_API_URL = "https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&q=%s"
)

type google struct {
	basicTranslator
}

func NewGoogle(cache SessionCache) Translator {
	return &google{
		basicTranslator{
			cache:  cache,
			agent:  DefaultAgent,
			engine: EngineGoogle,
		},
	}
}

func (b *google) Session() (*Session, error) {
	return &Session{}, nil
}

type GoogleResult struct {
	data interface{}
}

func (gr *GoogleResult) Data() []interface{} {
	if d, ok := gr.data.([]interface{}); ok {
		return d
	}
	return nil
}

func (gr *GoogleResult) String() (res string) {
	defer func() {
		if r := recover(); r != nil {
			res = ""
		}
	}()
	if d, ok := gr.data.([]interface{}); ok {
		if d1, ok := d[0].([]interface{}); ok {
			if len(d1) > 1 {
				for i := range d1[:len(d1)-1] {
					if d2, ok := d1[i].([]interface{}); ok {
						res += d2[0].(string)
					}
				}
			}
			return
		}
	}
	return
}

func (g *google) Translate(srcLang, targetLang, text string) (r Result, err error) {
	resp, err := g.get(fmt.Sprintf(GOOGLE_API_URL, srcLang, targetLang, url.QueryEscape(text)))
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		gr := []interface{}{}
		err = json.Unmarshal(b, &gr)
		if err == nil {
			r = &GoogleResult{data: gr}
		}
	}
	return
}
