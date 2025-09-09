package adif

import (
	"encoding/json"
	"strings"
	"testing"
)

func BenchmarkADIWrite(b *testing.B) {
	qsoListNative := loadTestData()
	b.ResetTimer()
	for b.Loop() {
		var sb strings.Builder
		w := NewADIRecordWriter(&sb)
		w.Write(qsoListNative)
		_ = sb.String()
	}
}

func BenchmarkADIJWrite(b *testing.B) {
	qsoListNative := loadTestData()
	b.ResetTimer()
	for b.Loop() {
		var sb strings.Builder
		w := NewJSONRecordWriter(&sb)
		w.Write(qsoListNative)
		_ = sb.String()
	}
}

// This benchmark works directly on JSON data, without using this library for writing the JSON ADIF data.
// It is meant as a reference point to compare the performance to a known standard (the go stdlib JSON parser).
func BenchmarkJSONWriteReference(b *testing.B) {
	jsonRecords := benchmarkFileAsJSON()
	document := adifDocument{}
	err := json.Unmarshal(jsonRecords, &document)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := json.Marshal(document)
		if err != nil {
			panic(err)
		}
	}
}
