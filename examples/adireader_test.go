package examples

import (
	"fmt"
	"io"
	"strings"

	"github.com/hamradiolog-net/adif"
	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

func ExampleADIFReader_Next() {
	var r = strings.NewReader("<PROGRAMID:7>MonoLog<EOH><CALL:5>W9PVA<EOR><CALL:5>K9CTS<EOR>")
	adiReader := adif.NewADIReader(r, true) // true means we'll skip the header (if there is one)

	for {
		qso, isHeader, bytesProcessed, err := adiReader.Next()
		if err == io.EOF {
			// io.EOF means there are no more records
			break
		}
		if err != nil {
			// this means that something went wrong.
			// see errors.go for errors specific to parsing.
			// other errors may be returned in addition to the ones listed in errors.go.
			panic(err)
		}

		fmt.Println(qso[adifield.CALL]) // n.b. the field keys must be UPPERCASE
		fmt.Println()

		fmt.Print(qso.String())
		if isHeader {
			fmt.Println(adif.TagEOH)
		} else {
			fmt.Println(adif.TagEOR)
		}
		fmt.Printf("Read %d bytes.\n\n", bytesProcessed)
	}

	// Output:
	// W9PVA
	//
	// <CALL:5>W9PVA<EOR>
	// Read 43 bytes.
	//
	// K9CTS
	//
	// <CALL:5>K9CTS<EOR>
	// Read 18 bytes.
}
