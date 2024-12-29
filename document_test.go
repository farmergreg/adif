package adif

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentString(t *testing.T) {
	// Arrange
	doc := &Document{}
	doc.ReadFrom(strings.NewReader("<PROGRAMID:7>MonoLog<EOH><CALL:5>W9PVA<EOR>"))

	// Act
	s := doc.String()

	// Assert
	assert.Equal(t, AdifHeaderPreamble+"<PROGRAMID:7>MonoLog<EOH>\n<CALL:5>W9PVA<EOR>\n", s)
}

func TestParseExportParseVerifySimple(t *testing.T) {
	// Arrange
	tests := []struct {
		name       string
		adifSource string
	}{
		{"Simple ADI", "<APP_MY_APP_0:0>\n<APP_MY_APP:1>x<EOR>"},
		{"Simple ADI with header", "<PROGRAMID:7>MonoLog<EOH>\n<APP_MY_APP_0:0><APP_MY_APP:1>x<EOR>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseExportParseVerifyHelper(t, tt.adifSource)
		})
	}
}

func TestParseExportParseVerifyFiles(t *testing.T) {
	t.Parallel()

	// Arrange
	files, err := testFileFS.ReadDir("testdata")
	assert.Nil(t, err)

	for _, file := range files {
		file := file
		t.Run(file.Name(), func(t *testing.T) {
			t.Parallel()

			// Act
			reader, err := testFileFS.Open("testdata/" + file.Name())
			assert.Nil(t, err)
			defer reader.Close()

			content, err := io.ReadAll(reader)
			assert.Nil(t, err)

			parseExportParseVerifyHelper(t, string(content))
		})
	}
}
func parseExportParseVerifyHelper(t *testing.T, adif string) {
	// Parse
	firstDoc := &Document{}
	firstDoc.ReadFrom(strings.NewReader(adif))

	// Export

	buf := &strings.Builder{}
	firstDoc.WriteTo(buf)

	// Parse
	secondDoc := &Document{}
	secondDoc.ReadFrom(strings.NewReader(buf.String()))

	// Verify
	if firstDoc.Header != nil {
		assert.Equal(t, len(firstDoc.Header.Fields), len(secondDoc.Header.Fields))
		for i := 0; i < len(firstDoc.Header.Fields); i++ {
			assert.Equal(t, firstDoc.Header.Fields[i], secondDoc.Header.Fields[i])
		}
	}
	assert.Equal(t, len(firstDoc.Records), len(secondDoc.Records))
	for r := range firstDoc.Records {
		assert.Equal(t, len(firstDoc.Records[r].Fields), len(secondDoc.Records[r].Fields))
		for i := 0; i < len(firstDoc.Records[r].Fields); i++ {
			assert.Equal(t, firstDoc.Records[r].Fields[i], secondDoc.Records[r].Fields[i])
		}
	}
}

func TestDocumentStringNilReceiver(t *testing.T) {
	// Arrange
	var doc *Document

	// Act
	s := doc.String()

	// Assert
	assert.Equal(t, "", s)
}
