package adif

import (
	"bytes"
	"encoding/json"
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

func BenchmarkADIJRead(b *testing.B) {
	var qsoList []Record
	json := benchmarkFileAsJSON()
	b.ResetTimer()
	for b.Loop() {
		qsoList = make([]Record, 0, 10000)
		p, err := NewJSONRecordReader(bytes.NewReader(json), false)
		if err != nil {
			b.Fatal(err)
		}
		if err != nil {
			b.Fatal(err)
		}
		for {
			q, err := p.Next()
			if err == io.EOF {
				break
			}
			qsoList = append(qsoList, q)
		}
	}

	_ = len(qsoList)
}

// This benchmark works directly on JSON data, without using this library for reading the JSON ADIF data.
// It is meant as a reference point to compare the performance to a known standard (the go stdlib JSON parser).
func BenchmarkJSONReadReference(b *testing.B) {
	jsonRecords := benchmarkFileAsJSON()

	b.ResetTimer()
	document := adifDocument{}
	for b.Loop() {
		err := json.Unmarshal(jsonRecords, &document)
		if err != nil {
			b.Fatal(err)
		}
	}
}
