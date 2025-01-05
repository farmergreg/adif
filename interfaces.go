package adif

// ADIFParser parses Amateur Data Interchange Format (ADIF) records sequentially.
type ADIFParser interface {

	// Parse reads and returns the next Record in the input.
	// It returns io.EOF when no more records are available.
	Parse() (record *Record, isHeader bool, bytesRead int64, err error)
}
