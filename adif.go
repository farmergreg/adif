// package adif provides
// 1) Types, Structs and Methods for managing ADIF QSOs.
// 2) ADIF Reader for ADI formatted data.
// 3) Export ADI formatted data.
package adif

import (
	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

// The version of this library.
// See https://semver.org/
//
//	MAJOR version      == Incompatible API changes.
//	MINOR version      == Added functionality in a backward compatible manner
//	PATCH version      == Backward compatible bug fixes.
//	PRERELEASE version == Pre-release version (optional). This should be empty ("") or start with a dash (e.g. "-rc1").
//
// Additional labels for pre-release and build metadata are available as extensions to the MAJOR.MINOR.PATCH format.
const (
	versionMajor      = 0
	versionMinor      = 0
	versionPatch      = 1
	versionPreRelease = "-alpha"
)

const (
	TagEOH = string("<" + adifield.EOH + ">") // TagEOH is the end of header ADIF tag: <EOH>
	TagEOR = string("<" + adifield.EOR + ">") // TagEOR is the end of record ADIF tag: <EOR>
)

// DocumentMaxSizeInBytes controls the maximum size of data read into an Document struct in bytes.
// This variable helps prevent memory exhaustion attacks.
// You may adjust this value to suit your needs.
//
// For large documents, consider using the ADI Reader to stream the records.
// The default limit is 256MB.
var DocumentMaxSizeInBytes int64 = 1024 * 1024 * 256
