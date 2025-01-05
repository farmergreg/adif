// package adif provides
// 1) Types, Structs and Methods for managing ADIF Records.
// 2) ADIF Reader for ADI formatted data.
// 3) Export ADI formatted data.
package adif

import (
	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

const (
	TagEOH = string("<" + adifield.EOH + ">") // TagEOH is the end of header ADIF tag: <EOH>
	TagEOR = string("<" + adifield.EOR + ">") // TagEOR is the end of record ADIF tag: <EOR>
)

const adiHeaderPreamble = "K9CTS AM✠DG ADIF Library\n"
