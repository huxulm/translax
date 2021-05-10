package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/jackdon/translax/pkg/translator"
)

var (
	ErrTextLengthOver = errors.New("文本长度超过5000")
)
var eng = flag.String("e", "google", "engine for translate, available values are google, bing and sougou, default use google.")
var from = flag.String("from", "", "source language")
var to = flag.String("to", "", "target language")
var text = flag.String("text", "", "text content")

func checkInputLength(text string) (bool, error) {
	if n := utf8.RuneCountInString(text); n > 5000 {
		return false, ErrTextLengthOver
	}
	return true, nil
}

func Run() bool {
	var args = os.Args
	if len(args) == 1 {
		return false
	}
	flag.Parse()
	if ok, err := checkInputLength(*text); !ok {
		fmt.Printf("%v, 忽略超出内容", err)
		*text = (*text)[:5000]
	}
	r, err := translator.Trans(translator.EngineName(*eng), *from, *to, *text)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
	return true
}
