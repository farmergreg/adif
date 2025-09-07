package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
)

func ExampleNewRecord() {
	record := NewRecord()
	record.Set(adifield.CALL, "W1AW")
	record.Set(adifield.BAND, "10m")
	record.Set(adifield.MODE, "SSB")
	record.Set(adifield.APP_+"K9CTS_TEST", "TEST")

	if record.Get(adifield.CALL) != "W1AW" {
		panic("Expected W1AW, got " + record.Get(adifield.CALL)) // n.b. the field keys must be UPPERCASE
	}

	fmt.Print(record.String())
	fmt.Println(TagEOR)

	// Output: <BAND:3>10m<MODE:3>SSB<CALL:4>W1AW<APP_K9CTS_TEST:4>TEST<EOR>
}

func ExampleRecord_ReadFrom() {
	adiStr := "<CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>"
	r := NewRecord()

	// ReadFrom reads exactly one ADIF record from the given reader.
	// It skips the header if present.
	_, err := r.ReadFrom(strings.NewReader(adiStr))
	if err != nil {
		// see errors.go for errors specific to parsing.
		// other errors may be returned in addition to the ones listed in errors.go.
		panic(err)
	}

	fmt.Println(r.Get(adifield.CALL)) // n.b. the field keys must be UPPERCASE
	fmt.Println()

	fmt.Print(r.String()) // n.b. the fields do not always appear in the same order
	fmt.Println(TagEOR)

	// Output:
	// W1AW
	//
	// <BAND:3>10m<MODE:3>SSB<CALL:4>W1AW<APP_K9CTS_TEST:4>TEST<EOR>
}

func ExampleRecord_WriteTo() {
	record := NewRecord()
	record.Set(adifield.CALL, "W1AW")
	record.Set(adifield.BAND, "10m")
	record.Set(adifield.MODE, "SSB")
	record.Set(adifield.APP_+"K9CTS_TEST", "TEST")

	sb := strings.Builder{}
	_, err := record.WriteTo(&sb)
	if err != nil {
		panic(err)
	}

	fmt.Print(sb.String())
	fmt.Println(TagEOR)

	// Output: <BAND:3>10m<MODE:3>SSB<CALL:4>W1AW<APP_K9CTS_TEST:4>TEST<EOR>
}
