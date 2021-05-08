package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/jackdon/translax/pkg/translator"
)

var eng = flag.String("engine", "google", "engine for translate, available values are google, bing and sougou, default use google.")
var from = flag.String("from", "", "source language")
var to = flag.String("to", "", "target language")
var text = flag.String("text", "", "text content")

func Run() bool {
	var args = os.Args
	if len(args) == 1 {
		return false
	}
	flag.Parse()
	r, err := translator.Trans(translator.EngineName(*eng), *from, *to, *text)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
	return true
}
