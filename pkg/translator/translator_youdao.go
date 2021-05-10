package translator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	YOUDAO_PAGE    = "https://fanyi.youdao.com/"
	YOUDAO_API_URL = "https://fanyi.youdao.com/translate_o?smartresult=dict&smartresult=rule"
)

type youdao struct {
	basicTranslator
}

func NewYoudao(cache SessionCache) Translator {
	return &youdao{
		basicTranslator{
			cache:  cache,
			agent:  DefaultAgent,
			engine: EngineYoudao,
		},
	}
}

func (t *youdao) Engine() EngineName {
	return EngineYoudao
}

func (t *youdao) Session() (*Session, error) {
	req, _ := http.NewRequest("GET", YOUDAO_PAGE, nil)
	req.Header.Set("User-Agent", t.agent)
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

// {"translateResult":[[{"tgt":"你好,杰克","src":"Hello, Jack"}]],"errorCode":0,"type":"en2zh-CHS"}
type YoudaoResult struct {
	TranslateResult [][]*struct {
		Tgt string `json:"tgt"`
		Src string `json:"src"`
	} `json:"translateResult,omitempty"`
	ErrorCode int    `json:"errorCode"` // 0,
	Type      string `json:"type"`      // "en2zh-CHS"
}

func (yr *YoudaoResult) String() string {
	if yr == nil || yr.ErrorCode != 0 || len(yr.TranslateResult) == 0 {
		return ""
	}
	var s []string
	for _, r := range yr.TranslateResult {
		for _, ri := range r {
			s = append(s, ri.Tgt)
		}
		s = append(s, "\n")
	}
	return strings.Join(s, "")
}

/*
i: Hello
from: AUTO
to: AUTO
smartresult: dict
client: fanyideskweb
salt: 16206193311948
//    1620619546497325070
sign: e08bf4fe5488a280bd62aa50ba4251bb
lts:   1620619331194
bv: 940ae85ba32dcefb13c38faaf66f115f
doctype: json
version: 2.1
keyfrom: fanyi.web
action: FY_BY_REALTlME
*/

/*
i: Hello, Jack
from: AUTO
to: AUTO
smartresult: dict
client: fanyideskweb
salt: 16206194077798
sign: 032c4dbce28ec457fb60bfcf7cf3fae9
lts: 1620619407779
bv: 940ae85ba32dcefb13c38faaf66f115f
doctype: json
version: 2.1
keyfrom: fanyi.web
action: FY_BY_REALTlME
*/

/*
{text: "{"translateResult":[[{"tgt":"你好,杰克","src":"Hello, Jack"}]],"errorCode":0,"type":"en2zh-CHS"}"}
*/
// 0. e = "xxx"
// 1. t = md5(agent)  => 940ae85ba32dcefb13c38faaf66f115f
// 2. r = Date.now().getTime()
// 3. i = r + rand.Intn(10)
// 4.
/*
	ts: r,
	bv: t,
	salt: i,
	sign: n.md5("fanyideskweb" + e + i + "Tbh5E8=q6U3EXe+&L[4c@")
*/
func sign(text, salt string) string {
	return md5V(fmt.Sprintf("fanyideskweb%s%sTbh5E8=q6U3EXe+&L[4c@", text, salt))
}

func (bn *youdao) prepareForm(sl, tl, text string) url.Values {
	rand.Seed(time.Now().Unix())
	lts := fmt.Sprintf("%d", time.Now().UnixNano())[:13]
	salt := lts + fmt.Sprintf("%d", rand.Intn(10))
	return url.Values{
		"i":           {text},
		"from":        {"AUTO"},
		"to":          {"AUTO"},
		"smartresult": {"dict"},
		"client":      {"fanyideskweb"},
		"salt":        {salt},
		"sign":        {sign(text, salt)},
		"bv":          {md5V("5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.146 Safari/537.36")},
		"lts":         {lts},
		"doctype":     {"json"},
		"version":     {"2.1"},
		"keyfrom":     {"fanyi.web"},
		"action":      {"FY_BY_REALTlME"},
	}
}
func (bn *youdao) Translate(srcLang, targetLang, text string) (r Result, err error) {
	sl, tl, err := bn.keepLang(srcLang, targetLang)
	if err != nil {
		return
	}
	resp, err := bn.postForm(YOUDAO_API_URL, bn.prepareForm(sl, tl, text))
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		br := YoudaoResult{}
		err = json.Unmarshal(b, &br)
		if err == nil {
			r = &br
		}
	}
	return
}
