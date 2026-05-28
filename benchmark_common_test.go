package adif

import (
	_ "embed"
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
				s := NewScanner(reader)
				for s.Scan() {
				}
				if err := s.Err(); err != nil {
					b.Fatal(err)
				}
				reader.Close()
			}
		})
	}
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
						b.Fatalf("expected error for input %q", td.input)
					}
				} else if v != td.want {
					b.Fatalf("got %d, want %d", v, td.want)
				}
			}
		})
	}
}

func loadTestData() []Record {
	var records []Record
	s := NewScanner(strings.NewReader(benchmarkFile))
	for s.Scan() {
		records = append(records, s.Record())
	}
	if err := s.Err(); err != nil {
		panic(err)
	}
	return records
}
