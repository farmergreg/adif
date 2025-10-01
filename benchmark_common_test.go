package adif

import (
	_ "embed"
	"errors"
	"io"
	"strings"
	"testing"
)

//go:embed testdata/N3FJP-AClogAdif.adi
var benchmarkFile string

func BenchmarkParseAllADIFiles(b *testing.B) {
	fs, err := testFileFS.ReadDir("testdata")
	if err != nil {
		b.Fatal(err)
	}

	for _, f := range fs {
		b.Run(f.Name(), func(b *testing.B) {
			for b.Loop() {
				reader, _ := testFileFS.Open("testdata/" + f.Name())
				p := NewADIDocumentReader(reader, true)
				_, _, err := p.Next()
				for err == nil {
					_, _, err = p.Next()
				}
				reader.Close()
			}
		})

	}
}

func loadTestData() []Record {
	var qsoListNative []Record
	p := NewADIDocumentReader(strings.NewReader(benchmarkFile), false)
	record, _, err := p.Next()
	for err == nil {
		qsoListNative = append(qsoListNative, record)
		record, _, err = p.Next()
	}
	if !errors.Is(err, io.EOF) {
		panic(err)
	}
	return qsoListNative
}

func BenchmarkInternalParseDataLength(b *testing.B) {
	testData := []struct {
		input []byte
		want  int
	}{
		{[]byte("1"), 1},
		{[]byte("12"), 12},
		{[]byte("123"), 123},
		{[]byte("XYZMORE"), -1},
	}

	for _, td := range testData {
		b.Run(string(td.input), func(b *testing.B) {
			for b.Loop() {
				v, err := parseDataLength(td.input)
				if td.want == -1 {
					if err == nil {
						b.Fatalf("Expected error for input '%s', got nil", td.input)
					}
				} else if v != td.want {
					b.Fatalf("Expected %d, got %d", td.want, v)
				}
			}
		})
	}
}
