package translator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	BING_PAGE_URL = "https://cn.bing.com/translator"
	BING_API_URL  = "https://cn.bing.com/ttranslatev3?isVertical=1&IG=A19BFC4FE4624DCA902152EC61EB67EE&IID=translator.5025.31"
)

type bing struct {
	basicTranslator
}

func NewBing(cache SessionCache) Translator {
	return &bing{
		basicTranslator{
			cache:  cache,
			agent:  DefaultAgent,
			engine: EngineBing,
		},
	}
}

func (b *bing) Engine() EngineName {
	return EngineBing
}

func (b *bing) Session() (*Session, error) {
	req, _ := http.NewRequest("GET", BING_PAGE_URL, nil)
	req.Header.Set("User-Agent", b.agent)
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		defer resp.Body.Close()
		return &Session{
			Cookies: resp.Cookies(),
		}, nil
	} else {
		return nil, err
	}
}

type BingResult struct {
	DetectedLanguage *struct {
		Language string  `json:"language"` // zh-Hans
		Score    float32 `json:"score"`    //1.0
	} `json:"detectedLanguage,omitempty"`
	Translations []struct {
		SentLen struct {
			SrcSentLen []int64 `json:"srcSentLen"`
		} `json:"sentLen"`
		Text string `json:"text"`
		To   string `json:"to"`
	} `json:"translations"`
	InputTransliteration *string `json:"inputTransliteration,omitempty"`
}

func (br *BingResult) String() string {
	if len(br.Translations) > 0 {
		return br.Translations[0].Text
	}
	return ""
}

func (bn *bing) Translate(srcLang, targetLang, text string) (r Result, err error) {
	sl, tl, err := bn.keepLang(srcLang, targetLang)
	if err != nil {
		return
	}
	resp, err := bn.postForm(BING_API_URL, url.Values{
		"fromLang": {sl},
		"to":       {tl},
		"text":     {text},
	})
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		br := []BingResult{}
		err = json.Unmarshal(b, &br)
		if err == nil && len(br) > 0 {
			r = &br[0]
		}
	}
	return
}
