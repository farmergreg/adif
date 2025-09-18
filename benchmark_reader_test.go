package adif

import (
	"io"
	"strings"
	"testing"
)

func BenchmarkADIRead(b *testing.B) {
	var qsoList []Record
	b.ResetTimer()
	for b.Loop() {
		qsoList = make([]Record, 0, 10000)
		p := NewADIRecordReader(strings.NewReader(benchmarkFile), false)
		for {
			q, err := p.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
			qsoList = append(qsoList, q)
		}
	}

	_ = len(qsoList)
}
