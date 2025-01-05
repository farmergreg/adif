package examples

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif"
)

func ExampleDocument_ReadFrom() {
	headerPreamble := "Example Alternate Header Preamble\n"
	adi := "<ADIF_VER:5>3.1.5<EOH><CALL:5>W9PVA<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR><CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>"
	doc := adif.NewDocument(headerPreamble)
	doc.ReadFrom(strings.NewReader(adi))

	fmt.Println(doc.String())
	// Output:
	// Example Alternate Header Preamble
	// <ADIF_VER:5>3.1.5<EOH>
	// <CALL:5>W9PVA<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
	// <CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
}

func ExampleDocument_WriteTo() {
	headerPreamble := "Example Alternate Header Preamble\n"
	adi := "<ADIF_VER:5>3.1.5<EOH><CALL:5>W9PVA<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR><CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>"
	doc := adif.NewDocument(headerPreamble)

	doc.ReadFrom(strings.NewReader(adi))

	sb := &strings.Builder{}
	doc.WriteTo(sb)
	fmt.Println(sb.String())
	// Output:
	// Example Alternate Header Preamble
	// <ADIF_VER:5>3.1.5<EOH>
	// <CALL:5>W9PVA<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
	// <CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
}
