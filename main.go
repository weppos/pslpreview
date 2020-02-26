package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/preview", PreviewServer)
	http.ListenAndServe(":"+port, nil)
}

func PreviewServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", r.URL)

	query := r.URL.Query()
	var value string
	pp := &PreviewParams{}

	if value = query.Get("r"); value == "" {
		fmt.Fprintf(w, "Parameter 'r' is missing")
		return
	}
	pp.setRules(value)

	if value = query.Get("h"); value == "" {
		fmt.Fprintf(w, "Parameter 'h' is missing")
		return
	}
	pp.setHosts(value)

	fmt.Fprintf(w, fmt.Sprintf("%v", pp))
}

type PreviewParams struct {
	Rules []string
	Hosts []string
}

func (p *PreviewParams) setRules(value string) error {
	p.Rules = strings.Split(value, ",")
	return nil
}

func (p *PreviewParams) setHosts(value string) error {
	p.Hosts = strings.Split(value, ",")
	return nil
}
