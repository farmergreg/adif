package adif

import (
	"strings"
	"testing"
)

func BenchmarkADIRead(b *testing.B) {
	var records []Record
	b.ResetTimer()
	for b.Loop() {
		records = make([]Record, 0, 10000)
		s := NewScanner(strings.NewReader(benchmarkFile))
		for s.Scan() {
			records = append(records, s.Record())
		}
		if err := s.Err(); err != nil {
			b.Fatal(err)
		}
	}
	_ = len(records)
}
