package adif

import (
	"strings"
	"testing"
)

func TestDataDocumentString(t *testing.T) {
	// Arrange
	parser := NewAdiParser()

	reader := strings.NewReader("<PROGRAMID:7>MonoLog<EOH><CALL:5>W9PVA<EOR><CALL:5>K9CTS<EOR>")
	parser.ReadFrom(reader)

	// Act
	s := parser.GetDocument().String()

	// Assert
	want := adiHeaderPreamble + "<PROGRAMID:7>MonoLog<EOH>\n<CALL:5>W9PVA<EOR>\n<CALL:5>K9CTS<EOR>\n"
	if s != want {
		t.Errorf("got %q, want %q", s, want)
	}
}
