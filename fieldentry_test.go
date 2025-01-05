package adif

import "testing"

func TestFieldEntry(t *testing.T) {
	fe := FieldEntry{
		Name: "CALL",
		Data: "KC9RYZ",
	}

	if fe.String() != "<CALL:6>KC9RYZ" {
		t.Errorf("FieldEntry.String() = %s", fe.String())
	}
}
