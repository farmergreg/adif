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

func BenchmarkADIWriteMode(b *testing.B) {
	records := loadTestData()

	modes := []struct {
		name string
		mode WriteMode
	}{
		{"Fast", ADIWriteModeFast},
		{"Pretty", ADIWriteModePretty},
	}

	for _, m := range modes {
		b.Run(m.name, func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				var sb strings.Builder
				w := NewWriter(&sb).SetWriteMode(m.mode)
				for _, r := range records {
					w.Write(r)
				}
				w.Flush()
				_ = sb.String()
			}
		})
	}
}
