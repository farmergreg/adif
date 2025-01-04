package examples

import (
	"fmt"
	"strings"

	"github.com/hamradiolog-net/adif"
	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

func ExampleNewRecord() {
	record := adif.NewRecord()
	record.Set(adifield.CALL, "W1AW")
	record.Set(adifield.BAND, "10m")
	record.Set(adifield.MODE, "SSB")
	record.Set(adifield.APP_+"K9CTS_TEST", "TEST")

	fmt.Println(record.String())
	// Output: <CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
}

func ExampleRecord_ReadFrom() {
	adi := "<CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>"
	r := adif.NewRecord()

	// ReadFrom reads exactly one ADIF record from the given reader.
	// It skips the header if present.
	_, err := r.ReadFrom(strings.NewReader(adi))
	if err != nil {
		// see errors.go for errors specific to parsing.
		// other errors may be returned in addition to the ones listed in errors.go.
		panic(err)
	}

	fmt.Println(r.String())
	// Output: <CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
}

func ExampleRecord_WriteTo() {
	record := adif.NewRecord()
	record.Set(adifield.CALL, "W1AW")
	record.Set(adifield.BAND, "10m")
	record.Set(adifield.MODE, "SSB")
	record.Set(adifield.APP_+"K9CTS_TEST", "TEST")

	sb := strings.Builder{}
	record.WriteTo(&sb)
	fmt.Println(sb.String())
	// Output: <CALL:4>W1AW<BAND:3>10m<MODE:3>SSB<APP_K9CTS_TEST:4>TEST<EOR>
}
