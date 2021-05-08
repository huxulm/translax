package main

import (
	"github.com/jackdon/translax/pkg/cli"
	"github.com/jackdon/translax/pkg/cli/interactive"
)

func main() {
	if ok := cli.Run(); !ok {
		interactive.Run()
	}
}
