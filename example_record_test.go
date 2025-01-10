package adif

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

func ExampleNewRecord() {
	record := NewRecord()
	record[adifield.CALL] = "W1AW"
	record[adifield.BAND] = "10m"
	record[adifield.MODE] = "SSB"
	record[adifield.APP_+"K9CTS_TEST"] = "TEST"

	if record[adifield.CALL] != "W1AW" {
		panic("Expected W1AW, got " + record[adifield.CALL]) // n.b. the field keys must be UPPERCASE
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

	fmt.Println(r[adifield.CALL]) // n.b. the field keys must be UPPERCASE
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
	record[adifield.CALL] = "W1AW"
	record[adifield.BAND] = "10m"
	record[adifield.MODE] = "SSB"
	record[adifield.APP_+"K9CTS_TEST"] = "TEST"

	sb := strings.Builder{}
	record.WriteTo(&sb)

	fmt.Print(sb.String())
	fmt.Println(TagEOR)

	// Output: <BAND:3>10m<MODE:3>SSB<CALL:4>W1AW<APP_K9CTS_TEST:4>TEST<EOR>
}
