package adif

import (
	"fmt"
	"strings"

	"github.com/farmergreg/spec/v6/adifield"
	"github.com/farmergreg/spec/v6/enum/band"
	"github.com/farmergreg/spec/v6/enum/mode"
)

func ExampleNewADIJReader() {
	jsonExample := `{
  "header": {
    "created_timestamp": "20250907 212700",
    "programid": "ExampleProgram",
    "programversion": "1.0"
  },
  "records": [
    {
      "band": "20m",
      "call": "K9CTS",
      "mode": "ssb",
      "qso_date": "20250907",
      "time_on": "2127"
    },
    {
      "band": "40m",
      "call": "W9PVA",
      "mode": "cw",
      "qso_date": "20250907",
      "time_on": "2130"
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
		record, isHeader, err := reader.Next()
		if err != nil {
			break // EOF or other error
		}

		fmt.Printf("Is Header: %v\n", isHeader)
		if isHeader {
			fmt.Printf("created_timestamp: %s\n", record.Get(adifield.CREATED_TIMESTAMP))
		} else {
			fmt.Printf("call: %s, band: %s, mode: %s\n", record.Get(adifield.CALL), record.Get(adifield.BAND), record.Get(adifield.MODE))
		}
		fmt.Println()
	}

	// Output:
	// Is Header: true
	// created_timestamp: 20250907 212700
	//
	// Is Header: false
	// call: K9CTS, band: 20m, mode: ssb
	//
	// Is Header: false
	// call: W9PVA, band: 40m, mode: cw
}

// ExampleNewADIJWriter demonstrates how to write ADIJ JSON document using NewADIJWriter.
func ExampleNewJSONWriter() {
	var sb strings.Builder
	writer := NewJSONWriter(&sb, "  ")

	hdr := NewRecord()
	hdr.Set(adifield.CREATED_TIMESTAMP, "20250907 212700")
	writer.Write(hdr, true)

	qso := NewRecord()
	qso.Set(adifield.CALL, "K9CTS")
	qso.Set(adifield.BAND, band.BAND_20M.String())
	qso.Set(adifield.MODE, mode.SSB.String())
	writer.Write(qso, false)

	writer.Close()
	fmt.Println(sb.String())

	// Output:
	// {
	//   "header": {
	//     "created_timestamp": "20250907 212700"
	//   },
	//   "records": [
	//     {
	//       "band": "20m",
	//       "call": "K9CTS",
	//       "mode": "ssb"
	//     }
	//   ]
	// }
}
