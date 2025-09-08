package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
)

func ExampleNewADIReader() {
	// Example ADI data
	adiData := `
<ADIF_VERS:5>3.1.0
<PROGRAMID:4>Test
<EOH>
<CALL:5>K9CTS<QSO_DATE:8>20230101<TIME_ON:4>1200<BAND:3>20M<MODE:3>SSB<eor>
<CALL:5>W9PVA<QSO_DATE:8>20230102<TIME_ON:4>1300<BAND:3>40M<MODE:2>CW<eor>
`

	reader := NewADIReader(strings.NewReader(adiData), true)
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
