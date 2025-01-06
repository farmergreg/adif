package adif

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	_ "embed"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
	/*
		eminlin "github.com/Eminlin/GoADIFLog"
		eminlinformat "github.com/Eminlin/GoADIFLog/format"
		matir "github.com/Matir/adifparser"
	*/)

func BenchmarkAllTestFiles(b *testing.B) {
	fs, err := testFileFS.ReadDir("testdata")
	if err != nil {
		b.Fatal(err)
	}

	for _, f := range fs {
		b.Run(f.Name(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				reader, _ := testFileFS.Open("testdata/" + f.Name())
				content, _ := io.ReadAll(reader)
				p := NewADIReader(strings.NewReader(string(content)), true)
				for {
					_, _, _, err := p.Next()
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

func loadTestData() []Record {
	var qsoListNative []Record
	p := NewADIReader(strings.NewReader(benchmarkFile), false)
	for {
		record, _, _, err := p.Next()
		if err == io.EOF {
			break
		}
		qsoListNative = append(qsoListNative, record)
	}
	return qsoListNative
}

func BenchmarkReadThisLibrary(b *testing.B) {
	var qsoList []Record
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qsoList = make([]Record, 0, 10000)
		p := NewADIReader(strings.NewReader(benchmarkFile), false)
		for {
			q, _, _, err := p.Next()
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

	var records []Record
	var readCountJSON int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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

/*
func BenchmarkReadMatir(b *testing.B) {
	var qsoList []matir.ADIFRecord
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qsoList = make([]matir.ADIFRecord, 0, 10000)
		r := matir.NewADIFReader(strings.NewReader(benchmarkFile))
		for {
			q, err := r.ReadRecord()
			if err != nil {
				break
			}
			qsoList = append(qsoList, q)
		}
	}

	_ = len(qsoList)
}

func BenchmarkReadEminlin(b *testing.B) {
	qsoListNative := loadTestData()

	file, err := os.CreateTemp("", "eminlin-test.adi")
	if err != nil {
		b.Fatal(err)
	}
	for _, qso := range qsoListNative {
		qso.WriteTo(file)
	}
	file.Close()
	realFile := file.Name() + ".adi"
	os.Rename(file.Name(), realFile)

	var log []eminlinformat.CQLog
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log, err = eminlin.Parse(realFile)
		if err != nil {
			b.Fatal(err)
		}
	}

	_ = len(log)
	os.Remove(realFile)
}
*/

func BenchmarkWriteThisLibrary(b *testing.B) {
	qsoListNative := loadTestData()
	var writeCountADI int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		writeCountADI = 0
		for _, qso := range qsoListNative {
			qso.WriteTo(&sb)
			writeCountADI++
		}
		_ = sb.String()
	}

	if len(qsoListNative) != writeCountADI {
		b.Errorf("Write count mismatch: ADI %d, expected %d", writeCountADI, len(qsoListNative))
	}
}

func BenchmarkWriteJSON(b *testing.B) {
	qsoListNative := loadTestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(qsoListNative)
		if err != nil {
			b.Fatal(err)
		}
		_ = string(data)
	}
}

/*
func BenchmarkWriteMatir(b *testing.B) {
	// Setup Matir test data
	var qsoListMatir []matir.ADIFRecord
	r := matir.NewADIFReader(strings.NewReader(benchmarkFile))
	for {
		qso, err := r.ReadRecord()
		if err != nil {
			break
		}
		qsoListMatir = append(qsoListMatir, qso)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		mw := matir.NewADIFWriter(&sb)
		for _, qso := range qsoListMatir {
			mw.WriteRecord(qso)
		}
		_ = sb.String()
	}
}
*/

/*
// Eminlin does not support writing adi files
func BenchmarkWriteEminlin(b *testing.B) {
}
*/
func BenchmarkRandomFieldAccess(b *testing.B) {
	// Load test data once before the benchmark
	qsoListNative := loadTestData()

	// Common fields to access randomly
	fields := []adifield.Field{"CALL", "BAND", "MODE", "QSO_DATE", "TIME_ON", "APP_K9CTS", "STATE"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Access a random record and field
		recordIdx := i % len(qsoListNative)
		fieldIdx := i % len(fields)

		record := qsoListNative[recordIdx]
		field := fields[fieldIdx]

		// Access the field and do something with it to prevent optimization
		_ = record.Get(field)
	}
}
