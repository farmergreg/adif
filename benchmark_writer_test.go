package adif

import (
	"strings"
	"testing"
)

func BenchmarkRecordWriteTo(b *testing.B) {
	records := loadTestData()
	b.ResetTimer()
	for b.Loop() {
		var sb strings.Builder
		for _, r := range records {
			r.WriteTo(&sb)
			sb.Reset()
		}
	}
}

func BenchmarkADIWrite(b *testing.B) {
	records := loadTestData()
	b.ResetTimer()
	for b.Loop() {
		var sb strings.Builder
		w := NewWriter(&sb)
		for _, r := range records {
			w.Write(r)
		}
		w.Flush()
		_ = sb.String()
	}
}
