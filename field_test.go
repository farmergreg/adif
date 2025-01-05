package adif

import "testing"

func TestField(t *testing.T) {
	fe := Field{
		Name: "CALL",
		Data: "KC9RYZ",
	}

	if fe.String() != "<CALL:6>KC9RYZ" {
		t.Errorf("Field.String() = %s", fe.String())
	}
}
