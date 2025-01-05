package adif

import (
	"fmt"
)

var _ fmt.Stringer = Field{}

// String returns the Field as an ADI formatted string.
func (f Field) String() string {
	return fmt.Sprintf("<%s:%d>%s", f.Name, len(f.Data), f.Data)
}
