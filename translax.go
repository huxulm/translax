package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jackdon/translax/pkg/translator"

	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "google", Description: "google translator"},
		{Text: "bing", Description: "bing translator"},
		{Text: "sougou", Description: "sougou translator"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func main() {
	defer handleExit()
	p := prompt.New(executor, completer)
	p.Run()
}

func executor(l string) {
	blocks := strings.Split(l, " ")
	var (
		r   string
		err error
	)
	switch blocks[0] {
	case "google":
		r, err = translator.Trans(blocks[0], blocks[1], blocks[2], blocks[3])
	default:
		r, err = translator.Trans(string(translator.EngineGoogle), "en", "zh", l)
	}
	if err == nil {
		fmt.Println(r)
	} else {
		fmt.Println(err)
	}
}
