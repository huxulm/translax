package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jackdon/translax/pkg/translator"

	"github.com/c-bata/go-prompt"
)

var commands = []prompt.Suggest{
	{Text: "google", Description: "google translator"},
	{Text: "bing", Description: "bing translator"},
	{Text: "sougou", Description: "sougou translator"},

	{Text: "exit", Description: "exit the program"},
}

var translatorOptions = []prompt.Suggest{
	{Text: "--from", Description: "原文本语言."},
	{Text: "--to", Description: "目标语言."},
	{Text: "--text", Description: "文本内容."},
}
var optionsHelp = []prompt.Suggest{}

func excludeOptions(args []string) ([]string, bool) {
	l := len(args)
	if l == 0 {
		return nil, false
	}
	cmd := args[0]
	filtered := make([]string, 0, l)

	var skipNextArg bool
	for i := 0; i < len(args); i++ {
		if skipNextArg {
			skipNextArg = false
			continue
		}

		if cmd == "logs" && args[i] == "-f" {
			continue
		}

		for _, s := range []string{
			"--from", "--to",
			"--text",
		} {
			if strings.HasPrefix(args[i], s) {
				if strings.Contains(args[i], "=") {
					// we can specify option value like '-o=json'
					skipNextArg = false
				} else {
					skipNextArg = true
				}
				continue
			}
		}
		if strings.HasPrefix(args[i], "-") {
			continue
		}

		filtered = append(filtered, args[i])
	}
	return filtered, skipNextArg
}

func optionCompleter(args []string, long bool) []prompt.Suggest {
	l := len(args)
	if l <= 1 {
		if long {
			return prompt.FilterHasPrefix(optionsHelp, "--", false)
		}
		return optionsHelp
	}

	var suggests []prompt.Suggest

	commandArgs, _ := excludeOptions(args)
	switch commandArgs[0] {
	case "google", "bing", "sougou":
		suggests = translatorOptions
	default:
		suggests = optionsHelp
	}
	if long {
		return prompt.FilterContains(
			prompt.FilterHasPrefix(suggests, "--", false),
			strings.TrimLeft(args[l-1], "--"),
			true,
		)
	}
	return prompt.FilterContains(suggests, strings.TrimLeft(args[l-1], "-"), true)
}

func argumentsCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}
	return []prompt.Suggest{}
}

func getPreviousOption(d prompt.Document) (cmd, option string, found bool) {
	args := strings.Split(d.TextBeforeCursor(), " ")
	l := len(args)
	if l >= 2 {
		option = args[l-2]
	}
	if strings.HasPrefix(option, "-") {
		return args[0], option, true
	}
	return "", "", false
}

func completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return commands
	}

	args := strings.Split(d.TextBeforeCursor(), " ")
	w := d.GetWordBeforeCursor()

	// If word before the cursor starts with "-", returns CLI flag options.
	if strings.HasPrefix(w, "-") {
		return optionCompleter(args, strings.HasPrefix(w, "--"))
	}

	// Return suggestions for option
	if suggests, found := completeOptionArguments(d); found {
		return suggests
	}
	commandArgs, skipNext := excludeOptions(args)
	if skipNext {
		return []prompt.Suggest{}
	}
	return argumentsCompleter(commandArgs)
}

func completeOptionArguments(d prompt.Document) ([]prompt.Suggest, bool) {
	cmd, option, found := getPreviousOption(d)
	if !found {
		return []prompt.Suggest{}, false
	}
	if option == "--from" || option == "--to" {
		return prompt.FilterHasPrefix(
				getLangListSuggestions(translator.LangMap),
				d.GetWordBeforeCursor(),
				true),
			true
	}
	switch cmd {
	// TODO
	}
	return []prompt.Suggest{}, false
}

func getLangListSuggestions(langlist map[string]string) []prompt.Suggest {
	var suggestions []prompt.Suggest
	for lang := range langlist {
		suggestions = append(suggestions, prompt.Suggest{
			Text: lang,
		})
	}
	return suggestions
}

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func main() {
	defer handleExit()
	p := prompt.New(executor, completer, prompt.OptionPrefix("> "))
	p.Run()
}

func executor(l string) {
	blocks := strings.Split(l, " ")
	var (
		r   string
		err error
	)
	var sl, tl, text string
	for i, block := range blocks {
		if strings.HasSuffix(block, "from") {
			sl = blocks[i+1]
		}
		if strings.HasSuffix(block, "to") {
			tl = blocks[i+1]
		}
	}
	if idx := strings.Index(l, "--text"); idx > 0 {
		text = strings.Trim(l[idx+6:], " ")
	}
	switch blocks[0] {
	case "google", "bing", "sougou":
		r, err = translator.Trans(translator.EngineName(blocks[0]), sl, tl, text)
	default:
		r, err = translator.Trans(translator.EngineGoogle, "auto", "zh", blocks[len(blocks)-1])
	}
	if err == nil {
		fmt.Println(r)
	} else {
		fmt.Println(err)
	}
}
