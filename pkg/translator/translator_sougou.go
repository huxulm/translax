package translator

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

const (
	SOUGOU_PAGE    = "https://fanyi.sogou.com/?keyword=&transfrom=auto&model=general"
	SOUGOU_API_URL = "https://fanyi.sogou.com/api/transpc/text/result"
)

type sougou struct {
	basicTranslator
}

func NewSougou(cache SessionCache) Translator {
	return &sougou{
		basicTranslator{
			cache:  cache,
			agent:  DefaultAgent,
			engine: EngineSougou,
		},
	}
}

func (b *sougou) Engine() EngineName {
	return EngineSougou
}

func (b *sougou) Session() (*Session, error) {
	req, _ := http.NewRequest("GET", SOUGOU_PAGE, nil)
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

type FanyiReqBody struct {
	Client string `json:"client"` //  "pc"
	Fr     string `json:"fr"`     //  "browser_pc"
	From   string `json:"from"`   //  "auto"
	NeedQc int    `json:"needQc"` //  1
	S      string `json:"s"`      // "669005a6f7dcc02aa83a798bf6d9cc23"
	Text   string `json:"text"`   //  "hello"
	To     string `json:"to"`     //  "zh-CHS"
	Uuid   string `json:"uuid"`   //  "cd3b0491-561f-4508-a580-890930827f34"
}

type SougouResult struct {
	Info   string      `json:"info"`
	Status interface{} `json:"status"`
	Data   struct {
		Translate *struct {
			Dit       string `json:"dit"`
			ErrorCode string `json:"errorCode"`
		} `json:"translate"`
		// Network *struct {
		// 	NetworkMean []string `json:"network_mean"`
		// } `json:"network"`
	} `json:"data"`
}

func (sr *SougouResult) String() string {
	return sr.Data.Translate.Dit
}

func prepareData(sl, tl, text string) *FanyiReqBody {
	b := &FanyiReqBody{
		Client: "pc",
		Fr:     "browser_pc",
		From:   sl,
		NeedQc: 1,
		Text:   text,
		To:     tl,
		Uuid:   "",
	}
	u, _ := uuid.NewUUID()
	b.Uuid = u.String()
	crypt(b)
	return b
}

func crypt(d *FanyiReqBody) {
	d.S = md5V(fmt.Sprintf("%s%s%s109984457", d.From, d.To, d.Text))
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (sg *sougou) Translate(srcLang, targetLang, text string) (r Result, err error) {
	sl, tl, err := sg.keepLang(srcLang, targetLang)
	if err != nil {
		return
	}
	var d []byte
	d, err = json.Marshal(prepareData(sl, tl, text))
	if err != nil {
		return
	}
	resp, err := sg.post(SOUGOU_API_URL, d)
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		r = new(SougouResult)
		err = json.Unmarshal(b, r)
		if err != nil {
			r = nil
		}
	}
	return
}
