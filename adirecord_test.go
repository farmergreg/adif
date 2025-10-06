package adif

import (
	"testing"

	"github.com/farmergreg/spec/v6/adifield"
	"github.com/farmergreg/spec/v6/enum/band"
)

func TestNewADIRecordWithCapacity(t *testing.T) {
	_ = newRecordWithCapacity(7)
}

func TestNewADIRecord(t *testing.T) {
	_ = NewRecord()
}

func TestADIRecordSet_AddField(t *testing.T) {
	r := NewRecord()
	r.Set(adifield.CALL, "K9CTS")

	if r.Get(adifield.CALL) != "K9CTS" {
		t.Errorf("Expected value 'K9CTS', got '%s'", r.Get(adifield.CALL))
	}
}

func TestADIRecordSet_RemoveField(t *testing.T) {
	r := NewRecord()
	r.Set(adifield.CALL, "K9CTS")
	r.Set(adifield.CALL, "")
}

func TestADIRecordAll(t *testing.T) {
	r := NewRecord()
	r.Set(adifield.CALL, "K9CTS")
	r.Set(adifield.BAND, band.BAND_20M.String())
	for k, v := range r.Fields() {
		if k == adifield.CALL && v != "K9CTS" {
			t.Errorf("Expected value 'K9CTS' for CALL, got '%s'", v)
		}
		break // Testing the iterator !yield condition
	}
}
