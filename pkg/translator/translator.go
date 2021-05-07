package translator

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var DefaultAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.146 Safari/537.36"

type EngineName string

const (
	EngineBaidu  = EngineName("baidu")
	EngineSougou = EngineName("sougou")
	EngineBing   = EngineName("bing")
	EngineGoogle = EngineName("google")
)

var langMap = map[string]string{
	"af":  "af",
	"sq":  "sq",
	"am":  "am",
	"ar":  "ar",
	"hy":  "hy",
	"as":  "as",
	"az":  "az",
	"bn":  "bn",
	"bs":  "bs",
	"bg":  "bg",
	"yue": "yue",
	"ca":  "ca",
	"zh":  "zh",
	"hr":  "hr",
	"cs":  "cs",
	"da":  "da",
	"prs": "prs",
	"nl":  "nl",
	"en":  "en",
	"et":  "et",
	"fj":  "fj",
	"fil": "fil",
	"fi":  "fi",
	"fr":  "fr",
	"de":  "de",
	"el":  "el",
	"gu":  "gu",
	"ht":  "ht",
	"he":  "he",
	"hi":  "hi",
	"mww": "mww",
	"hu":  "hu",
	"is":  "is",
	"id":  "id",
	"iu":  "iu",
	"ga":  "ga",
	"it":  "it",
	"ja":  "ja",
	"kn":  "kn",
	"kk":  "kk",
	"km":  "km",
	"tlh": "tlh",
	"ko":  "ko",
	"ku":  "ku",
	"kmr": "kmr",
	"lo":  "lo",
	"lv":  "lv",
	"lt":  "lt",
	"mg":  "mg",
	"ms":  "ms",
	"ml":  "ml",
	"mt":  "mt",
	"mi":  "mi",
	"mr":  "mr",
	"my":  "my",
	"ne":  "ne",
	"nb":  "nb",
	"or":  "or",
	"ps":  "ps",
	"fa":  "fa",
	"pl":  "pl",
	"pt":  "pt",
	"pa":  "pa",
	"otq": "otq",
	"ro":  "ro",
	"ru":  "ru",
	"sm":  "sm",
	"sr":  "sr",
	"sk":  "sk",
	"sl":  "sl",
	"es":  "es",
	"sw":  "sw",
	"sv":  "sv",
	"ty":  "ty",
	"ta":  "ta",
	"te":  "te",
	"th":  "th",
	"ti":  "ti",
	"to":  "to",
	"tr":  "tr",
	"uk":  "uk",
	"ur":  "ur",
	"vi":  "vi",
	"cy":  "cy",
	"yua": "yua",
}

type Result interface {
	fmt.Stringer
}

type Session struct {
	ExprAt  int64          `yaml:"expr_at"`
	Cookies []*http.Cookie `yaml:"cookies"`
}

type SessionCache interface {
	Persist(engine EngineName, session *Session) error
	GetSession(engine EngineName) (*Session, error)
	GetTranslatorByEngineName(engine EngineName) Translator
}

type defaultSessionCache struct {
	SessionCache
	memSession map[EngineName]*Session
}

func (c *defaultSessionCache) Load() error {
	dir, err := getDir()
	if err != nil {
		return err
	}
	for e := range ENGINES {
		d, err := ioutil.ReadFile(filepath.Join(dir, string(e)+".yaml"))
		if err != nil {
			return err
		}
		s := new(Session)
		if err := yaml.Unmarshal(d, s); err == nil {
			c.memSession[e] = s
		}
	}
	return nil
}

func getDir() (dir string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	dir = filepath.Join(homeDir, ".config", "translaX")
	return
}
func (c *defaultSessionCache) Persist(engine EngineName, session *Session) error {
	if session == nil {
		return errors.New("session can not be nil.")
	}
	dir, err := getDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("缓存目录创建失败: %v", err)
	}
	if d, err := yaml.Marshal(session); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filepath.Join(dir, string(engine)+".yaml"), d, os.ModePerm)
	}
}

func (c *defaultSessionCache) GetSession(engine EngineName) (*Session, error) {
	if s, ok := c.memSession[engine]; ok {
		return s, nil
	} else {
		t := c.GetTranslatorByEngineName(engine)
		if t == nil {
			return nil, errors.New("no translator found")
		}
		s, err := t.Session()
		if err == nil {
			c.memSession[engine] = s
			c.Persist(engine, s)
		}
		return s, err
	}
}

func (c *defaultSessionCache) GetTranslatorByEngineName(engine EngineName) Translator {
	return ENGINES[engine]
}

type Translator interface {
	Engine() EngineName
	Session() (*Session, error)
	Translate(srcLang, targetLang, text string) (Result, error)
	postForm(url string, data url.Values) (*http.Response, error)
	post(url string, data []byte) (*http.Response, error)
}

type basicTranslator struct {
	Translator
	engine EngineName
	agent  string
	cache  SessionCache
}

func (b *basicTranslator) Engine() EngineName {
	return b.engine
}

func (b *basicTranslator) postForm(url string, data url.Values) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest("POST", url, strings.NewReader(data.Encode())); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s, err := b.cache.GetSession(b.Engine())
	if err != nil {
		return
	}
	for _, c := range s.Cookies {
		req.Header.Add("Cookie", c.Raw)
	}
	req.Header.Set("User-Agent", b.agent)
	return http.DefaultClient.Do(req)
}

func (b *basicTranslator) post(url string, data []byte) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest("POST", url, bytes.NewBuffer(data)); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	s, err := b.cache.GetSession(b.Engine())
	if err != nil {
		return
	}
	for _, c := range s.Cookies {
		req.Header.Add("Cookie", c.Raw)
	}
	req.Header.Set("User-Agent", b.agent)
	return http.DefaultClient.Do(req)
}

func (b *basicTranslator) get(url string) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", b.agent)
	return http.DefaultClient.Do(req)
}

func (b *basicTranslator) keepLang(srcLang, targetLang string) (sl, tl string, err error) {
	if sl, okSl := langMap[strings.ToLower(srcLang)]; okSl {
		if tl, okTl := langMap[strings.ToLower(targetLang)]; okTl {
			if b.Engine() == EngineBing {
				if sl == "zh" {
					sl = "zh-Hans"
				}
				if tl == "zh" {
					tl = "zh-Hans"
				}
			}
			if b.Engine() == EngineSougou {
				if sl == "zh" {
					sl = "zh-CHS"
				}
				if tl == "zh" {
					tl = "zh-CHS"
				}
			}
			return sl, tl, nil
		}
	}
	return "", "", errors.New("not supported language code.")
}

var ENGINES = map[EngineName]Translator{}

func RegisterTranslator(translator Translator) {
	ENGINES[translator.Engine()] = translator
}

var defaultCache = &defaultSessionCache{
	memSession: make(map[EngineName]*Session),
}

func init() {
	RegisterTranslator(NewSougou(defaultCache))
	RegisterTranslator(NewBing(defaultCache))
	RegisterTranslator(NewGoogle(defaultCache))
	// after register all translator
	defaultCache.Load()
}

func Trans(engine, from, to, text string) (string, error) {
	switch engine {
	case string(EngineGoogle):
		r, err := ENGINES[EngineGoogle].Translate(from, to, text)
		return fmt.Sprintf("%v", r), err
	case string(EngineBing):
		r, err := ENGINES[EngineBing].Translate(from, to, text)
		return fmt.Sprintf("%v", r), err
	case string(EngineSougou):
		r, err := ENGINES[EngineSougou].Translate(from, to, text)
		return fmt.Sprintf("%v", r), err
	default:
		return "", nil
	}
}
