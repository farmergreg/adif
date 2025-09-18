package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/spec/v6/adifield"
	"github.com/hamradiolog-net/spec/v6/enum/band"
	"github.com/hamradiolog-net/spec/v6/enum/mode"
)

func ExampleNewADIRecordReader() {
	// Example ADI data
	adiData := `
<ADIF_VERS:5>3.1.0
<PROGRAMID:4>Test
<EOH>
<CALL:5>K9CTS<QSO_DATE:8>20230101<TIME_ON:4>1200<BAND:3>20M<MODE:3>SSB<eor>
<CALL:5>W9PVA<QSO_DATE:8>20230102<TIME_ON:4>1300<BAND:3>40M<MODE:2>CW<eor>
`

	reader := NewADIRecordReader(strings.NewReader(adiData), true)
	for {
		record, err := reader.Next()
		if err != nil {
			break // EOF or other error
		}
		fmt.Printf("Call: %s, Date: %s, Time: %s, Band: %s, Mode: %s\n",
			record.Get(adifield.CALL),
			record.Get(adifield.QSO_DATE),
			record.Get(adifield.TIME_ON),
			record.Get(adifield.BAND),
			record.Get(adifield.MODE))
	}
	// Output:
	// Call: K9CTS, Date: 20230101, Time: 1200, Band: 20M, Mode: SSB
	// Call: W9PVA, Date: 20230102, Time: 1300, Band: 40M, Mode: CW
}

// ExampleNewADIRecordWriter demonstrates how to write an ADI document using NewADIRecordWriter.
func ExampleNewADIRecordWriter() {
	var sb strings.Builder
	writer := NewADIRecordWriter(&sb)

	hdr := NewRecord()
	hdr.SetIsHeader(true)
	hdr.Set(adifield.CREATED_TIMESTAMP, "20250907 212700")

	qso := NewRecord()
	qso.Set(adifield.CALL, "K9CTS")
	qso.Set(adifield.BAND, band.BAND_20M.String())
	qso.Set(adifield.MODE, mode.SSB.String())
	qso.Set(adifield.New("APP_Example"), "Example")

	records := []Record{hdr, qso}
	for _, r := range records {
		writer.Write(r)
	}

	fmt.Println(sb.String())

	// Output:
	// AM✠DG
	// K9CTS High Performance ADIF Processing Library
	//    https://github.com/hamradiolog-net/adif
	//
	// <created_timestamp:15>20250907 212700<eoh>
	// <band:3>20m<mode:3>ssb<call:5>K9CTS<app_example:7>Example<eor>
}
