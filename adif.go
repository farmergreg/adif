// Package adif implements a high performance ADIF library for Go.
// It provides types, structs and methods for managing ADIF Records.
// Idiomatic interfaces for reading and writing ADI formatted data
// make integration with other Go libraries simple.
package adif

import (
	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

const (
	TagEOH = string("<" + adifield.EOH + ">") // TagEOH is the end of header ADIF tag: <EOH>
	TagEOR = string("<" + adifield.EOR + ">") // TagEOR is the end of record ADIF tag: <EOR>
)

const adiHeaderPreamble = "K9CTS AM✠DG ADIF Library\n"
