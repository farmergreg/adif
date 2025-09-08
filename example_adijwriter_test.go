package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
)

// ExampleNewADIJWriter demonstrates how to write ADIJ JSON document using NewADIJWriter.
func ExampleNewADIJWriter() {
	var sb strings.Builder
	writer := NewADIJWriter(&sb)

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
	// {
	//   "HEADER": {
	//     "CREATED_TIMESTAMP": "20250907 212700"
	//   },
	//   "RECORDS": [
	//     {
	//       "BAND": "20M",
	//       "CALL": "K9CTS",
	//       "MODE": "SSB"
	//     }
	//   ]
	// }
}
