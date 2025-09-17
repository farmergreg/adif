package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/spec/v6/adifield"
	"github.com/hamradiolog-net/spec/v6/enum/band"
	"github.com/hamradiolog-net/spec/v6/enum/mode"
)

// ExampleNewJSONRecordReader demonstrates how to read ADIJ JSON document using NewJSONRecordReader.
func ExampleNewJSONRecordReader() {
	jsonExample := `{
  "header": {
    "created_timestamp": "20250907 212700",
    "programid": "ExampleProgram",
    "programversion": "1.0"
  },
  "records": [
    {
      "band": "20M",
      "call": "K9CTS",
      "mode": "SSB",
      "qso_date": "20250907",
      "time_on": "2127"
    },
    {
      "band": "40M",
      "call": "W9PVA",
      "mode": "CW",
      "qso_date": "20250907",
      "time_on": "2130"
    }
  ]
}`

	// Create a reader from the ADIJ data
	reader, err := NewJSONRecordReader(strings.NewReader(jsonExample), false)
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

// ExampleNewJSONWriter demonstrates how to write ADIJ JSON document using NewJSONWriter.
func ExampleNewJSONRecordWriter() {
	var sb strings.Builder
	writer := NewJSONRecordWriter(&sb)

	hdr := NewRecord()
	hdr.SetIsHeader(true)
	hdr.Set(adifield.CREATED_TIMESTAMP, "20250907 212700")

	qso := NewRecord()
	qso.Set(adifield.CALL, "K9CTS")
	qso.Set(adifield.BAND, band.BAND_20M.String())
	qso.Set(adifield.MODE, mode.SSB.String())

	writer.Write([]Record{hdr, qso})

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
