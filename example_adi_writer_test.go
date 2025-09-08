package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
)

// ExampleNewADIWriter demonstrates how to write an ADI document using NewADIWriter.
func ExampleNewADIWriter() {
	var sb strings.Builder
	writer := NewADIWriter(&sb)

	hdr := NewADIRecord()
	hdr.SetIsHeader(true)
	hdr.Set(adifield.CREATED_TIMESTAMP, "20250907 212700")

	qso := NewADIRecord()
	qso.Set(adifield.BAND, "20M")
	qso.Set(adifield.CALL, "K9CTS")
	qso.Set(adifield.MODE, "SSB")

	writer.Write([]ADIFRecord{hdr, qso})

	fmt.Println(sb.String())

	// Output:
	// AM✠DG
	// K9CTS High Performance ADIF Processing Library
	// https://github.com/hamradiolog-net/adif-parser
	//
	// <CREATED_TIMESTAMP:15>20250907 212700<EOH>
	// <BAND:3>20M<MODE:3>SSB<CALL:5>K9CTS<EOR>
}
