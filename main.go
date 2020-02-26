package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

var (
	// Version is replaced at compilation time
	Version string
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

	w.Header().Set("X-Version", Version)

	query := r.URL.Query()
	var value string
	pp := &PreviewParams{}

	//if value = query.Get("r"); value == "" {
	//	fmt.Fprintf(w, "Parameter 'r' is missing")
	//	return
	//}
	//pp.setRules(value)

	if value = query.Get("h"); value == "" {
		fmt.Fprintf(w, "Parameter 'h' is missing")
		return
	}
	pp.setHosts(value)

	prs := preview(pp)

	fmt.Fprintf(w, fmt.Sprintf("%v", prs))
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

type PreviewResults struct {
	Results []PreviewResult
}

type PreviewResult struct {
	Host   string
	Domain PreviewDomain
	Error  error
}

type PreviewDomain struct {
	ETLD        string
	ETLDPlusOne string
	Rule        string
}

func newPreviewDomainFromPublicSuffix(d *publicsuffix.DomainName) *PreviewDomain {
	if d == nil {
		return nil
	}

	pd := &PreviewDomain{}
	pd.ETLD = d.TLD
	pd.ETLDPlusOne = d.SLD + "." + d.TLD

	rule := d.Rule.Value
	if d.Rule.Type == publicsuffix.WildcardType {
		rule = "*." + rule
	}
	if d.Rule.Type == publicsuffix.ExceptionType {
		rule = "!" + rule
	}
	pd.Rule = rule

	return pd
}

func preview(pp *PreviewParams) *PreviewResults {
	prs := &PreviewResults{}

	for _, h := range pp.Hosts {
		pr := PreviewResult{Host: h}

		d, err := publicsuffix.Parse(pr.Host)
		pr.Domain = *newPreviewDomainFromPublicSuffix(d)
		pr.Error = err

		prs.Results = append(prs.Results, pr)
	}

	return prs
}
