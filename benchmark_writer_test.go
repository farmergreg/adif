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
		w := NewADIRecordWriter(&sb)
		for _, r := range qsoListNative {
			w.Write(r)
		}
		_ = sb.String()
	}
}
