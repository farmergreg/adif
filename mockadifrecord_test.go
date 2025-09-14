package adif

import "github.com/hamradiolog-net/spec/v6/adifield"

var _ Record = (*mockADIFRecord)(nil)

// mockADIFRecord is a fake ADIFRecord used for testing.
type mockADIFRecord struct{}

func (r *mockADIFRecord) IsHeader() bool {
	return false
}

func (r *mockADIFRecord) SetIsHeader(isHeader bool) {}

func (r *mockADIFRecord) Get(field adifield.Field) string {
	return ""
}

func (r *mockADIFRecord) Set(field adifield.Field, value string) {}

func (r *mockADIFRecord) All() func(func(adifield.Field, string) bool) {
	return func(yield func(adifield.Field, string) bool) {
		yield(adifield.COMMENT, "")
	}
}

func (r *mockADIFRecord) Count() int {
	return 1
}

func (r *mockADIFRecord) GetDataType(field adifield.Field) rune {
	return 'S'
}

func (r *mockADIFRecord) SetDataType(field adifield.Field, dataType rune) {}
