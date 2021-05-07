package translator

import (
	"fmt"
	"testing"
)

func TestBingTranslate(t *testing.T) {

	tl, ok := ENGINES[EngineBing]
	if ok {
		r, err := tl.Translate("en", "zh", "Have you eaten?")
		if err == nil {
			fmt.Println(r.String())
		}
	}
}
