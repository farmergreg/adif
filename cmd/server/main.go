package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
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
	contentType     = "Content-Type"
	contentTypeJSON = "application/json"
	contentTypeADI  = "text/x-adif-adi"
)

var indexTemplate = template.Must(template.New("index").Parse(indexHTML))

func main() {
	addr := flag.String("addr", "localhost:8080", "server address to listen on")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/", handleIndex)
	mux.HandleFunc("POST /api/v1/", handleConversion)

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
	w.WriteHeader(http.StatusOK)
	indexTemplate.Execute(w, map[string]string{
		"ContentTypeADI":  contentTypeADI,
		"ContentTypeJSON": contentTypeJSON,
	})
}

func handleConversion(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get(contentType) {
	case contentTypeJSON: // JSON to ADIF
		var doc adif.Document
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			http.Error(w, "invalid json input", http.StatusBadRequest)
			return
		}

		w.Header().Set(contentType, contentTypeADI)
		w.WriteHeader(http.StatusOK)
		if _, err := doc.WriteTo(w); err != nil {
			http.Error(w, "unable to write adi output", http.StatusInternalServerError)
			return
		}
	case contentTypeADI: // ADIF to JSON
		doc := adif.NewDocument()
		if _, err := doc.ReadFrom(r.Body); err != nil {
			http.Error(w, "unable to read adi input", http.StatusBadRequest)
			return
		}

		w.Header().Set(contentType, contentTypeJSON)
		w.WriteHeader(http.StatusOK)

		jsonData, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			http.Error(w, "unable to create json output", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)

	default:
		http.Error(w,
			fmt.Sprintf("Unsupported %s. Use %s or %s", contentType, contentTypeADI, contentTypeJSON),
			http.StatusUnsupportedMediaType)
	}
}
