package adif

import "github.com/hamradiolog-net/spec/v6/adifield"

var _ ADIFRecord = (*mockADIFRecord)(nil)

// mockADIFRecord is a fake ADIFRecord used for testing.
type mockADIFRecord struct{}

func (r *mockADIFRecord) IsHeader() bool {
	return false
}

func (r *mockADIFRecord) SetIsHeader(isHeader bool) {}

func (r *mockADIFRecord) Get(field adifield.ADIField) string {
	return ""
}

func (r *mockADIFRecord) Set(field adifield.ADIField, value string) {}

func (r *mockADIFRecord) All() func(func(adifield.ADIField, string) bool) {
	return func(yield func(adifield.ADIField, string) bool) {
		yield(adifield.COMMENT, "")
	}
}

func (r *mockADIFRecord) Count() int {
	return 1
}
