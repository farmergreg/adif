package examples

import (
	"fmt"
	"io"
	"strings"

	"github.com/hamradiolog-net/adif"
)

func ExampleADIFParser_Parse() {
	var r = strings.NewReader("<PROGRAMID:7>MonoLog<EOH><CALL:5>W9PVA<EOR><CALL:5>K9CTS<EOR>")
	parser := adif.NewADIParser(r, true)

	for {
		record, _, err := parser.Parse()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		fmt.Println(record.String())
	}

	// Output:
	// <CALL:5>W9PVA<EOR>
	// <CALL:5>K9CTS<EOR>
}
