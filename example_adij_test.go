package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
	"github.com/hamradiolog-net/adif-spec/v6/enum/band"
	"github.com/hamradiolog-net/adif-spec/v6/enum/mode"
)

// ExampleNewADIJReader demonstrates how to read ADIJ JSON document using NewADIJReader.
func ExampleNewADIJReader() {
	jsonExample := `{
  "HEADER": {
    "CREATED_TIMESTAMP": "20250907 212700",
    "PROGRAMID": "ExampleProgram",
    "PROGRAMVERSION": "1.0"
  },
  "RECORDS": [
    {
      "BAND": "20M",
      "CALL": "K9CTS",
      "MODE": "SSB",
      "QSO_DATE": "20250907",
      "TIME_ON": "2127"
    },
    {
      "BAND": "40M",
      "CALL": "W9PVA",
      "MODE": "CW",
      "QSO_DATE": "20250907",
      "TIME_ON": "2130"
    }
  ]
}`

	// Create a reader from the ADIJ data
	reader, err := NewADIJReader(strings.NewReader(jsonExample), false)
	if err != nil {
		fmt.Printf("Error creating reader: %v\n", err)
		return
	}

	for {
		record, err := reader.Next()
		if err != nil {
			break // EOF or other error
		}

		fmt.Printf("Is Header: %v\n", record.IsHeader())
		if record.IsHeader() {
			fmt.Printf("Header CREATED_TIMESTAMP: %s\n", record.Get(adifield.CREATED_TIMESTAMP))
		} else {
			fmt.Printf("QSO CALL: %s, BAND: %s, MODE: %s\n", record.Get(adifield.CALL), record.Get(adifield.BAND), record.Get(adifield.MODE))
		}
		fmt.Println()
	}

	// Output:
	// Is Header: true
	// Header CREATED_TIMESTAMP: 20250907 212700
	//
	// Is Header: false
	// QSO CALL: K9CTS, BAND: 20M, MODE: SSB
	//
	// Is Header: false
	// QSO CALL: W9PVA, BAND: 40M, MODE: CW
}

// ExampleNewADIJWriter demonstrates how to write ADIJ JSON document using NewADIJWriter.
func ExampleNewADIJWriter() {
	var sb strings.Builder
	writer := NewADIJWriter(&sb)

	hdr := NewADIRecord()
	hdr.SetIsHeader(true)
	hdr.Set(adifield.CREATED_TIMESTAMP, "20250907 212700")

	qso := NewADIRecord()
	qso.Set(adifield.CALL, "K9CTS")
	qso.Set(adifield.BAND, band.Band20m.String())
	qso.Set(adifield.MODE, mode.SSB.String())

	writer.Write([]ADIFRecord{hdr, qso})

	fmt.Println(sb.String())

	// Output:
	// {
	//   "HEADER": {
	//     "CREATED_TIMESTAMP": "20250907 212700"
	//   },
	//   "RECORDS": [
	//     {
	//       "BAND": "20m",
	//       "CALL": "K9CTS",
	//       "MODE": "SSB"
	//     }
	//   ]
	// }
}
