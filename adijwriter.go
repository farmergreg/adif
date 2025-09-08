package adif

import (
	"encoding/json"
	"io"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
)

var _ ADIFRecordWriter = (*adijWriter)(nil)

// adijWriter implements ADIFRecordWriter for writing ADIF records in ADIJ format.
type adijWriter struct {
	w io.Writer
}

func NewADIJWriter(w io.Writer) *adijWriter {
	return &adijWriter{w: w}
}

func (aw *adijWriter) Write(records []ADIFRecord) (n int64, err error) {
	doc := &ADIJDocument{}
	if len(records) > 0 && records[0].IsHeader() {
		doc.Header = adijRecordToMap(records[0])
		records = records[1:]
	}

	for _, record := range records {
		doc.Records = append(doc.Records, adijRecordToMap(record))
	}

	encoder := json.NewEncoder(aw.w)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(doc)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func adijRecordToMap(record ADIFRecord) map[adifield.ADIField]string {
	result := make(map[adifield.ADIField]string)
	for _, field := range record.Fields() {
		value := record.Get(field)
		if value != "" {
			result[field] = value
		}
	}
	return result
}
