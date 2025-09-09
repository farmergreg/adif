package adif

import (
	"testing"

	"github.com/hamradiolog-net/spec/v6/adifield"
)

func TestNewADIRecordWithCapacity(t *testing.T) {
	_ = NewADIRecordWithCapacity(10)
}

func TestNewADIRecord(t *testing.T) {
	_ = NewADIRecord()
}

func TestADIRecordSet_AddField(t *testing.T) {
	r := NewADIRecord()
	r.Set(adifield.CALL, "K9CTS")
	if len(r.Fields()) != 1 {
		t.Errorf("Expected field count '1', got '%d'", len(r.Fields()))
	}
	if r.Get(adifield.CALL) != "K9CTS" {
		t.Errorf("Expected value 'K9CTS', got '%s'", r.Get(adifield.CALL))
	}
}

func TestADIRecordSet_RemoveField(t *testing.T) {
	r := NewADIRecord()
	r.Set(adifield.CALL, "K9CTS")
	r.Set(adifield.CALL, "")
	if len(r.Fields()) != 0 {
		t.Errorf("Expected field count '0', got '%d'", len(r.Fields()))
	}
}

func TestADIRecordSet_IsHeader(t *testing.T) {
	r := NewADIRecord()
	if r.IsHeader() {
		t.Errorf("Expected IsHeader false, got true")
	}
	r.SetIsHeader(true)
	if !r.IsHeader() {
		t.Errorf("Expected IsHeader true, got false")
	}
}
