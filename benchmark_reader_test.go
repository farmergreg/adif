package adif

import (
	"errors"
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
		q, _, err := p.Next()
		for err == nil {
			qsoList = append(qsoList, q)
			q, _, err = p.Next()
		}
		if !errors.Is(err, io.EOF) {
			b.Fatal(err)
		}
	}

	_ = len(qsoList)
}
