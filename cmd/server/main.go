package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/hamradiolog-net/adif"
)

//go:embed static
var staticFiles embed.FS

//go:embed template/index.html
var indexHTML string

const (
	contentType           = "Content-Type"
	contentTypeJSON       = "application/x-adif-json"
	contentTypeADI        = "application/x-adif-adi"
	contentTypeXML        = "application/x-adif-xml"
	contentTypeAutoDetect = "application/x-adif-auto-detect"
)

var indexTemplate = template.Must(template.New("index").Parse(indexHTML))

func main() {
	addr := flag.String("addr", "localhost:8080", "server address to listen on")
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("GET /", RequestLogger(http.HandlerFunc(handleIndex)))
	mux.Handle("POST /", RequestLogger(http.HandlerFunc(handleConversion)))

	// Serve static files (HTMX and PicoCSS)
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create sub-filesystem: %v", err)
	}

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	log.Printf("Starting server on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, struct {
		ContentTypeADI        string
		ContentTypeXML        string
		ContentTypeJSON       string
		ContentTypeAutoDetect string
	}{
		ContentTypeADI:        contentTypeADI,
		ContentTypeXML:        contentTypeXML,
		ContentTypeJSON:       contentTypeJSON,
		ContentTypeAutoDetect: contentTypeAutoDetect,
	})
}

func handleConversion(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get(contentType)
	data := r.Body

	if contentType == contentTypeAutoDetect {
		input, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "unable to read http request body", http.StatusBadRequest)
			return
		}

		// Naively assume that an ADIF file would not be audacious enough to begin with a comment that contains json
		input = bytes.TrimLeft(input, " \t\r\n")
		if input[0] == '{' || input[0] == '[' {
			contentType = contentTypeJSON
		} else {
			contentType = contentTypeADI
		}
		data = io.NopCloser(bytes.NewReader(input))
	}

	switch contentType {
	case contentTypeJSON: // JSON to ADIF
		var doc adif.Document
		if err := json.NewDecoder(data).Decode(&doc); err != nil {
			http.Error(w, "invalid json input", http.StatusBadRequest)
			return
		}

		w.Header().Set(contentType, contentTypeADI)
		if _, err := doc.WriteTo(w); err != nil {
			http.Error(w, "unable to write adi output", http.StatusInternalServerError)
			return
		}
	case contentTypeADI: // ADIF to JSON
		doc := adif.NewDocument()
		if _, err := doc.ReadFrom(data); err != nil {
			http.Error(w, "unable to read adi input", http.StatusBadRequest)
			return
		}

		w.Header().Set(contentType, contentTypeJSON)

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(doc); err != nil {
			http.Error(w, "unable to create json output", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w,
			fmt.Sprintf("Unsupported %s. Use %s or %s", contentType, contentTypeADI, contentTypeJSON),
			http.StatusUnsupportedMediaType)
	}

}
