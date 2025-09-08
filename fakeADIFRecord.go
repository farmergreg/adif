package adif

import "github.com/hamradiolog-net/adif-spec/v6/adifield"

var _ ADIFRecord = (*fakeADIFRecord)(nil)

// fakeADIFRecord is a fake ADIFRecord used for testing.
type fakeADIFRecord struct{}

func (r *fakeADIFRecord) IsHeader() bool {
	return false
}

func (r *fakeADIFRecord) SetIsHeader(isHeader bool) {}

func (r *fakeADIFRecord) Get(field adifield.ADIField) string {
	return ""
}

func (r *fakeADIFRecord) Set(field adifield.ADIField, value string) {}

func (r *fakeADIFRecord) Fields() []adifield.ADIField {
	return []adifield.ADIField{adifield.COMMENT}
}
