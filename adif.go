// package adif provides
// Types, Structs and Methods for managing ADIF Records.
// Interfaces that accept io.Reader and io.Writer for reading and writing ADI formatted data.
package adif

import (
	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

const (
	TagEOH = string("<" + adifield.EOH + ">") // TagEOH is the end of header ADIF tag: <EOH>
	TagEOR = string("<" + adifield.EOR + ">") // TagEOR is the end of record ADIF tag: <EOR>
)

const adiHeaderPreamble = "K9CTS AM✠DG ADIF Library\n"
