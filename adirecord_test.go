package adif

import (
	"testing"

	"github.com/hamradiolog-net/spec/v6/adifield"
	"github.com/hamradiolog-net/spec/v6/aditype"
	"github.com/hamradiolog-net/spec/v6/enum/band"
)

func TestNewADIRecordWithCapacity(t *testing.T) {
	_ = newRecordWithCapacity(10)
}

func TestNewADIRecord(t *testing.T) {
	_ = NewRecord()
}

func TestADIRecordSet_AddField(t *testing.T) {
	r := NewRecord()
	r.Set(adifield.CALL, "K9CTS")
	if r.Count() != 1 {
		t.Errorf("Expected field count '1', got '%d'", r.Count())
	}

	if r.Get(adifield.CALL) != "K9CTS" {
		t.Errorf("Expected value 'K9CTS', got '%s'", r.Get(adifield.CALL))
	}
}

func TestADIRecordSet_RemoveField(t *testing.T) {
	r := NewRecord()
	r.Set(adifield.CALL, "K9CTS")
	r.Set(adifield.CALL, "")
	if r.Count() != 0 {
		t.Errorf("Expected field count '0', got '%d'", r.Count())
	}
}

func TestADIRecordSet_IsHeader(t *testing.T) {
	r := NewRecord()
	if r.IsHeader() {
		t.Errorf("Expected IsHeader false, got true")
	}
	r.SetIsHeader(true)
	if !r.IsHeader() {
		t.Errorf("Expected IsHeader true, got false")
	}
}

func TestADIRecordAll(t *testing.T) {
	r := NewRecord()
	r.Set(adifield.CALL, "K9CTS")
	r.Set(adifield.BAND, band.BAND_20M.String())
	for k, v := range r.All() {
		if k == adifield.CALL && v != "K9CTS" {
			t.Errorf("Expected value 'K9CTS' for CALL, got '%s'", v)
		}
		break // Testing the iterator !yield condition
	}
}

func TestSetDataTypeOnKnownFieldNotPossible(t *testing.T) {
	r := NewRecord()
	f := adifield.CALL
	r.SetDataType(f, aditype.NewDataTypeIndicator('M'))
	if r.GetDataType(f) != aditype.NewDataTypeIndicator('S') {
		t.Errorf("Expected data type 'S', got '%s'", r.GetDataType(f))
	}
}

func TestGetDataTypeForAPP_(t *testing.T) {
	r := NewRecord()
	f := adifield.New("APP_TEST")
	if r.GetDataType(f) != aditype.NewDataTypeIndicator('M') {
		t.Errorf("Expected data type 'M', got '%s'", r.GetDataType(f))
	}
}

func TestSetGetDataTypeForAPP_(t *testing.T) {
	r := NewRecord()
	f := adifield.New("APP_TEST")
	r.SetDataType(f, aditype.NewDataTypeIndicator('S'))
	if r.GetDataType(f) != aditype.NewDataTypeIndicator('S') {
		t.Errorf("Expected data type 'S', got '%s'", r.GetDataType(f))
	}
}

func TestGetDataTypeForKnownField(t *testing.T) {
	r := NewRecord()
	f := adifield.BAND
	if r.GetDataType(f) != aditype.NewDataTypeIndicator('E') {
		t.Errorf("Expected data type 'E', got '%s'", r.GetDataType(f))
	}
}

func TestGetDataTypeForUnKnownField(t *testing.T) {
	r := NewRecord()
	f := adifield.New("__fake_future_adif_specification_field__")
	if r.GetDataType(f) != aditype.NewDataTypeIndicator(0) {
		t.Errorf("Expected data type '', got '%s'", r.GetDataType(f))
	}
}
