package doctrans

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"text/scanner"

	"github.com/stretchr/testify/assert"
)

func TestExtractPDFVersion(t *testing.T) {
	in := "/home/bx/Public/docs/tech_specs/pdf_reference_1-7.pdf"
	f, err := os.OpenFile(in, os.O_RDONLY, os.ModePerm)
	defer f.Close()
	if err == nil {
		// read first line
		firstLineEnd := seekFirstLineEnd(f)
		fmt.Println(firstLineEnd)
		var l = make([]byte, firstLineEnd+1)
		f.ReadAt(l, 0)
		fmt.Println(string(l))
	} else {
		t.Fatal(err)
	}
}

func seekFirstLineEnd(f *os.File) (l int64) {
	l = 0
	var c = make([]byte, 1)
	f.Seek(0, io.SeekStart)
	for {
		if _, err := f.Seek(l, 1); err != nil {
			return
		} else {
			if n, err := f.ReadAt(c, l); err == nil {
				if string(c[:n]) == "\n" {
					return
				}
				l++
				continue
			}
		}
	}
}

func TestExtractPDFVersion2(t *testing.T) {
	in := "/home/bx/Public/docs/tech_specs/pdf_reference_1-7.pdf"
	var s scanner.Scanner
	r, err := os.Open(in)
	pdfVersionPrefix := "%PDF-"
	var version string
	assert.NoError(t, err)
	s.Init(r)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		line := s.TokenText()
		if strings.HasPrefix(line, pdfVersionPrefix) {
			t.Logf("find version: %s", line)
			version = strings.Replace(line, pdfVersionPrefix, "", 1)
			break
		}
	}
	assert.EqualValues(t, "1.7", version)
}
