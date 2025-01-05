package adif

import (
	"fmt"
)

var _ fmt.Stringer = FieldEntry{}

// String returns the FieldEntry as an ADI formatted string.
func (f FieldEntry) String() string {
	return fmt.Sprintf("<%s:%d>%s", f.Name, len(f.Data), f.Data)
}
