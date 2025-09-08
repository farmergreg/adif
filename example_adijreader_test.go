package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
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
      "CALL": "N1ABC",
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

	// Read all records and collect them
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
	// QSO CALL: N1ABC, BAND: 40M, MODE: CW
}
