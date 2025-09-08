package adif

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	_ "embed"
)

func BenchmarkAllTestFiles(b *testing.B) {
	fs, err := testFileFS.ReadDir("testdata")
	if err != nil {
		b.Fatal(err)
	}

	for _, f := range fs {
		b.Run(f.Name(), func(b *testing.B) {
			for b.Loop() {
				reader, _ := testFileFS.Open("testdata/" + f.Name())
				p := NewADIReader(reader, true)
				for {
					_, _, err := p.Next()
					if err == io.EOF {
						break
					}
				}
				reader.Close()
			}
		})

	}
}

//go:embed testdata/N3FJP-AClogAdif.adi
var benchmarkFile string

func loadTestData() []ADIFRecord {
	var qsoListNative []ADIFRecord
	p := NewADIReader(strings.NewReader(benchmarkFile), false)
	for {
		record, _, err := p.Next()
		if err == io.EOF {
			break
		}
		qsoListNative = append(qsoListNative, record)
	}
	return qsoListNative
}

func BenchmarkReadThisLibrary(b *testing.B) {
	var qsoList []ADIFRecord
	b.ResetTimer()
	for b.Loop() {
		qsoList = make([]ADIFRecord, 0, 10000)
		p := NewADIReader(strings.NewReader(benchmarkFile), false)
		for {
			q, _, err := p.Next()
			if err == io.EOF {
				break
			}
			qsoList = append(qsoList, q)
		}
	}

	_ = len(qsoList)
}

func BenchmarkReadJSON(b *testing.B) {
	qsoListNative := loadTestData()
	jsonData, err := json.Marshal(qsoListNative)
	if err != nil {
		b.Fatal(err)
	}
	jsonString := string(jsonData)

	var records []adiRecord
	var readCountJSON int
	b.ResetTimer()
	for b.Loop() {
		records = records[:0]
		// convoluted, but this is to match the other benchmarks which also work with string input...
		// in reality, this does not affect the speed of this benchmark...
		err := json.Unmarshal([]byte(jsonString), &records)
		if err != nil {
			b.Fatal(err)
		}
		readCountJSON = len(records)
	}

	if len(qsoListNative) != readCountJSON {
		b.Errorf("Read count mismatch: JSON %d, expected %d", readCountJSON, len(qsoListNative))
	}
}

func BenchmarkWriteThisLibrary(b *testing.B) {
	qsoListNative := loadTestData()

	b.ResetTimer()
	for b.Loop() {
		var sb strings.Builder
		w := NewADIWriter(&sb)
		w.Write(qsoListNative)
		_ = sb.String()
	}

}
