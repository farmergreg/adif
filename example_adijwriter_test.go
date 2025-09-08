package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
)

func ExampleNewADIJWriter() {

	var sb strings.Builder
	writer := NewADIJWriter(&sb)

	qso := NewADIRecord()
	qso.SetIsHeader(false)
	qso.Set(adifield.BAND, "20M")
	qso.Set(adifield.CALL, "K9CTS")
	qso.Set(adifield.MODE, "SSB")

	writer.Write([]ADIFRecord{qso})

	fmt.Println(sb.String())

	// Output:
	// {
	//   "RECORDS": [
	//     {
	//       "BAND": "20M",
	//       "CALL": "K9CTS",
	//       "MODE": "SSB"
	//     }
	//   ]
	// }
}
