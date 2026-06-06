package adif_test

import (
	"fmt"
	"strings"

	adif "github.com/farmergreg/adif/v5"
	"github.com/farmergreg/spec/v6/adifield"
	"github.com/farmergreg/spec/v6/enum/band"
	"github.com/farmergreg/spec/v6/enum/mode"
)

// ExampleScanner demonstrates streaming QSO records from an ADI source.
func ExampleScanner() {
	adiData := `
<ADIF_VER:5>3.1.0
<PROGRAMID:4>Test
<EOH>
<CALL:5>K9CTS<QSO_DATE:8>20230101<TIME_ON:4>1200<BAND:3>20m<MODE:3>ssb<eor>
<CALL:5>W9PVA<QSO_DATE:8>20230102<TIME_ON:4>1300<BAND:3>40m<MODE:2>cw<eor>
`

	s := adif.NewScanner(strings.NewReader(adiData))
	for s.Scan() {
		if s.IsHeader() {
			continue
		}
		r := s.Record()
		fmt.Printf("Call: %s, Date: %s, Time: %s, Band: %s, Mode: %s\n",
			r[adifield.CALL], r[adifield.QSO_DATE], r[adifield.TIME_ON],
			r[adifield.BAND], r[adifield.MODE])
	}
	if err := s.Err(); err != nil {
		panic(err)
	}

	// Output:
	// Call: K9CTS, Date: 20230101, Time: 1200, Band: 20m, Mode: ssb
	// Call: W9PVA, Date: 20230102, Time: 1300, Band: 40m, Mode: cw
}

// ExampleWriter demonstrates writing an ADI document record by record.
func ExampleWriter() {
	var sb strings.Builder
	w := adif.NewWriter(&sb).SetWriteMode(adif.WriteModePretty)

	hdr := adif.NewRecord()
	hdr[adifield.CREATED_TIMESTAMP] = "20250907 212700"
	w.WriteHeader(hdr)

	qso := adif.NewRecord()
	qso[adifield.CALL] = "K9CTS"
	qso[adifield.BAND] = band.BAND_20M.String()
	qso[adifield.MODE] = mode.SSB.String()
	qso[adifield.New("APP_Example")] = "Example"

	// Pretty is the default; this is just for demonstration of SetWriteMode.
	w.SetWriteMode(adif.WriteModePretty)
	w.Write(qso)

	if err := w.Flush(); err != nil {
		panic(err)
	}

	fmt.Println(sb.String())

	// Output:
	// AM✠DG
	// K9CTS High Performance ADIF Processing Library
	//    https://github.com/farmergreg/adif
	//
	// <CREATED_TIMESTAMP:15>20250907 212700<EOH>
	// <BAND:3>20M<MODE:3>SSB<CALL:5>K9CTS<APP_EXAMPLE:7>Example<EOR>
}

// ExampleDocument demonstrates loading a complete ADI file into memory.
func ExampleDocument() {
	adi := "<PROGRAMID:4>Test<EOH><CALL:5>W9PVA<BAND:3>20m<EOR><CALL:5>K9CTS<BAND:3>40m<EOR>"

	d := adif.NewDocument()
	if _, err := d.ReadFrom(strings.NewReader(adi)); err != nil {
		panic(err)
	}

	fmt.Printf("Program: %s\n", d.Header[adifield.PROGRAMID])
	for _, r := range d.Records {
		fmt.Printf("Call: %s, Band: %s\n", r[adifield.CALL], r[adifield.BAND])
	}

	// Output:
	// Program: Test
	// Call: W9PVA, Band: 20m
	// Call: K9CTS, Band: 40m
}
