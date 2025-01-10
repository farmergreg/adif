package adif

import (
	"fmt"
	"io"
	"strings"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

func ExampleNewADIReader() {
	var r = strings.NewReader("<PROGRAMID:7>MonoLog<EOH><CALL:5>W9PVA<EOR><CALL:5>K9CTS<EOR>")
	adiReader := NewADIReader(r, true) // true means we'll skip the header (if there is one)

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

		fmt.Println(qso.Get(adifield.CALL))
		fmt.Println()

		fmt.Print(qso.String()) // n.b. the fields do not always appear in the same order
		if isHeader {
			fmt.Println(TagEOH)
		} else {
			fmt.Println(TagEOR)
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
