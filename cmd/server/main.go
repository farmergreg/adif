package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/hamradiolog-net/adif"
)

const (
	contentType     = "Content-Type"
	contentTypeJSON = "application/json"
	contentTypeADI  = "text/x-adif-adi"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "server address to listen on")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/", handleIndex)
	mux.HandleFunc("POST /api/v1/", handleConversion)

	log.Printf("Starting server on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil, `<h1>ADI <==> JSON Converter</h1>
<div>To convert ADI to JSON:</div>
<ul>
	<li>POST an ADI document to this endpoint with Content-Type %s</li>
	<li>Receive JSON output</li>
</ul>
<div>To convert JSON to ADI:</div>
<ul>
	<li>POST JSON to this endpoint with Content-Type %s</li>
	<li>Receive ADI output</li>
</ul>`,
		contentTypeADI,
		contentTypeJSON))
}

func handleConversion(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get(contentType) {
	case contentTypeJSON: // JSON to ADIF
		var doc adif.Document
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			http.Error(w, "Invalid JSON input.", http.StatusBadRequest)
			return
		}

		w.Header().Set(contentType, contentTypeADI)
		w.WriteHeader(http.StatusOK)
		if _, err := doc.WriteTo(w); err != nil {
			http.Error(w, "Failed to write ADI output.", http.StatusInternalServerError)
			return
		}
	case contentTypeADI: // ADIF to JSON
		doc := adif.NewDocument()
		if _, err := doc.ReadFrom(r.Body); err != nil {
			http.Error(w, "ADI input is invalid.", http.StatusBadRequest)
			return
		}

		w.Header().Set(contentType, contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(doc); err != nil {
			http.Error(w, "Failed to write JSON output.", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, fmt.Sprintf("Unsupported %s. Use %s or %s", contentType, contentTypeADI, contentTypeJSON), http.StatusUnsupportedMediaType)
	}
}
