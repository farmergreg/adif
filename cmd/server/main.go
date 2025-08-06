package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/hamradiolog-net/adif"
)

//go:embed static
var staticFiles embed.FS

//go:embed template/index.html
var indexHTML string

const (
	contentType     = "Content-Type"
	contentTypeJSON = "application/x-adif-json"
	contentTypeADI  = "application/x-adif-adi"
	contentTypeXML  = "application/x-adif-xml"
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
	indexTemplate.Execute(w, nil)
}

func handleConversion(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.Header.Get(contentType), "multipart/form-data;") {
		handleMulipartFormData(w, r)
	} else {
		handleSimpleApi(w, r)
	}
}

func handleSimpleApi(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get(contentType)
	data := r.Body

	switch contentType {
	case contentTypeADI: // ADIF to JSON
		convertAdi(w, data, func(w http.ResponseWriter) {
			w.Header().Set(contentType, contentTypeJSON)
		})

	case contentTypeJSON: // JSON to ADIF
		convertAdij(w, data, func(w http.ResponseWriter) {
			w.Header().Set(contentType, contentTypeADI)
		})

	default:
		http.Error(w,
			fmt.Sprintf("Unsupported %s. Use %s or %s", contentType, contentTypeADI, contentTypeJSON),
			http.StatusUnsupportedMediaType)
	}
}

func handleMulipartFormData(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 256) // 256KiB
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading input file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(path.Ext(handler.Filename))
	filename := handler.Filename[:len(handler.Filename)-len(ext)]
	switch ext {
	case ".adi":
		convertAdi(w, file, func(w http.ResponseWriter) {
			w.Header().Set("Content-Disposition", "attachment; filename="+filename+".adij")
		})

	case ".adij":
		convertAdij(w, file, func(w http.ResponseWriter) {
			w.Header().Set("Content-Disposition", "attachment; filename="+filename+".adi")
		})

	default:
		http.Error(w, "files must have .adi or .adij file extension", http.StatusUnsupportedMediaType)
	}
}

func convertAdi(w http.ResponseWriter, file io.Reader, beforeWriteCallback func(w http.ResponseWriter)) {
	doc := adif.NewDocument()
	if _, err := doc.ReadFrom(file); err != nil {
		http.Error(w, "unable to read adi input", http.StatusBadRequest)
		return
	}

	beforeWriteCallback(w)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(doc); err != nil {
		http.Error(w, "unable to create json output", http.StatusInternalServerError)
		return
	}
}

func convertAdij(w http.ResponseWriter, file io.Reader, beforeWriteCallback func(w http.ResponseWriter)) {
	var doc adif.Document
	if err := json.NewDecoder(file).Decode(&doc); err != nil {
		http.Error(w, "invalid json input", http.StatusBadRequest)
		return
	}

	beforeWriteCallback(w)
	if _, err := doc.WriteTo(w); err != nil {
		http.Error(w, "unable to write adi output", http.StatusInternalServerError)
		return
	}
}
