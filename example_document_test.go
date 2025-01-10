package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

func ExampleDocument_ReadFrom() {
	adi := "<ADIF_VER:5>3.1.5<EOH><CALL:5>W9PVA<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR><CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>"
	doc := NewDocument()
	doc.ReadFrom(strings.NewReader(adi))

	fmt.Println(doc.Records[0].Get(adifield.CALL)) // n.b. the field keys must be UPPERCASE
	fmt.Println()

	fmt.Println(doc.String())

	// Output:
	// W9PVA
	//
	// K9CTS AM✠DG ADIF Library
	// <ADIF_VER:5>3.1.5<EOH>
	// <CALL:5>W9PVA<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
	// <CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
}

func ExampleDocument_WriteTo() {
	adi := "<ADIF_VER:5>3.1.5<EOH><CALL:5>W9PVA<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR><CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>"
	doc := NewDocument()

	doc.ReadFrom(strings.NewReader(adi))

	sb := &strings.Builder{}
	doc.WriteTo(sb)

	fmt.Println(doc.Records[0].Get(adifield.CALL)) // n.b. the field keys must be UPPERCASE
	fmt.Println()

	fmt.Println(sb.String())

	// Output:
	// W9PVA
	//
	// K9CTS AM✠DG ADIF Library
	// <ADIF_VER:5>3.1.5<EOH>
	// <CALL:5>W9PVA<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
	// <CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
}
