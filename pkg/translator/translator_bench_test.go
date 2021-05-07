package translator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkSougouTranslator(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tl, ok := ENGINES[EngineSougou]
		if ok {
			_, err := tl.Translate("zh-CHS", "en", "你好！")
			assert.NoError(b, err)
		}
	}
}

func BenchmarkBingTranslator(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tl, ok := ENGINES[EngineBing]
		if ok {
			_, err := tl.Translate("zh-Hans", "en", "你好！")
			assert.NoError(b, err)
		}
	}
}
