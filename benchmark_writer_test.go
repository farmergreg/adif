package adif

import (
	"strings"
	"testing"
)

func BenchmarkADIWrite(b *testing.B) {
	qsoListNative := loadTestData()
	b.ResetTimer()
	for b.Loop() {
		var sb strings.Builder
		w := NewADIDocumentWriter(&sb)
		for _, r := range qsoListNative {
			w.WriteRecord(r)
		}
		_ = sb.String()
	}
}
