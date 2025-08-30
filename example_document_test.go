package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif-spec/v2/src/pkg/adifield"
)

func ExampleDocument_ReadFrom() {
	adi := "<ADIF_VER:5>3.1.5<EOH><CALL:5>W9PVA<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR><CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>"
	doc := NewDocument()
	doc.ReadFrom(strings.NewReader(adi))

	fmt.Println(doc.Records[0][adifield.CALL]) // n.b. the field keys must be UPPERCASE
	fmt.Println()

	fmt.Println(doc.String())

	// Output:
	// W9PVA
	//
	//                     AM✠DG
	// K9CTS High Performance ADIF Processing Library
	//     https://github.com/hamradiolog-net/adif
	//
	// <ADIF_VER:5>3.1.5<EOH>
	// <BAND:3>10m<MODE:3>SSB<CALL:5>W9PVA<APP_K9CTS_TEST:4>TEST<EOR>
	// <BAND:3>10m<MODE:3>SSB<CALL:4>W1AW<APP_K9CTS_TEST:4>TEST<EOR>
}

func ExampleDocument_WriteTo() {
	adi := "<ADIF_VER:5>3.1.5<EOH><CALL:5>W9PVA<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR><CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>"
	doc := NewDocumentWithOptions(2, "Example ADI Header Preamble\n")

	doc.ReadFrom(strings.NewReader(adi))

	sb := &strings.Builder{}
	doc.WriteTo(sb)

	fmt.Println(doc.Records[0][adifield.CALL]) // n.b. the field keys must be UPPERCASE
	fmt.Println()

	fmt.Println(sb.String()) // n.b. the fields do not always appear in the same order

	// Output:
	// W9PVA
	//
	// Example ADI Header Preamble
	// <ADIF_VER:5>3.1.5<EOH>
	// <BAND:3>10m<MODE:3>SSB<CALL:5>W9PVA<APP_K9CTS_TEST:4>TEST<EOR>
	// <BAND:3>10m<MODE:3>SSB<CALL:4>W1AW<APP_K9CTS_TEST:4>TEST<EOR>
}
