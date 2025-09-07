package adif

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	_ "embed"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
)

func BenchmarkAllTestFiles(b *testing.B) {
	fs, err := testFileFS.ReadDir("testdata")
	if err != nil {
		b.Fatal(err)
	}

	for _, f := range fs {
		b.Run(f.Name(), func(b *testing.B) {
			for b.Loop() {
				reader, _ := testFileFS.Open("testdata/" + f.Name())
				p := NewADIReader(reader, true)
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
	for b.Loop() {
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
	for b.Loop() {
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

func BenchmarkWriteThisLibrary(b *testing.B) {
	qsoListNative := loadTestData()
	var writeCountADI int

	b.ResetTimer()
	for b.Loop() {
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

func BenchmarkLoTWOneRecord(b *testing.B) {
	const oneRecord = `<APP_LoTW_OWNCALL:5>K9CTS
<STATION_CALLSIGN:5>K9CTS
<MY_DXCC:3>291
<MY_COUNTRY:24>UNITED STATES OF AMERICA
<APP_LoTW_MY_DXCC_ENTITY_STATUS:7>Current
<MY_GRIDSQUARE:6>EN34QU
<MY_STATE:2>WI // Wisconsin
<MY_CNTY:9>WI,PIERCE // Pierce
<MY_CQ_ZONE:2>04
<MY_ITU_ZONE:2>07
<CALL:5>N5ILQ
<BAND:3>20M
<FREQ:8>14.06100
<MODE:2>CW
<APP_LoTW_MODEGROUP:2>CW
<QSO_DATE:8>20220602
<APP_LoTW_RXQSO:19>2022-06-02 18:24:14 // QSO record inserted/modified at LoTW
<TIME_ON:6>182054
<APP_LoTW_QSO_TIMESTAMP:20>2022-06-02T18:20:54Z // QSO Date & Time; ISO-8601
<QSL_RCVD:1>Y
<QSLRDATE:8>20220602
<APP_LoTW_RXQSL:19>2022-06-02 23:31:22 // QSL record matched/modified at LoTW
<DXCC:3>291
<COUNTRY:24>UNITED STATES OF AMERICA
<APP_LoTW_DXCC_ENTITY_STATUS:7>Current
<PFX:2>N5
<APP_LoTW_2xQSL:1>Y
<GRIDSQUARE:4>EM15
<STATE:2>OK // Oklahoma
<CNTY:11>OK,OKLAHOMA // Oklahoma
<CQZ:2>04
<ITUZ:2>07
<eor>`

	record := NewRecord()

	b.ResetTimer()
	for b.Loop() {
		record.ReadFrom(strings.NewReader(oneRecord))
		_ = record.r[adifield.CALL]
	}
}
