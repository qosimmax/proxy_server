package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Content types
var CONTENT_TYPES = []string{"html", "text", "json", "xml"}

//Proxy struct
type Proxy struct {
	targetURL string // target host url
	oldStr    string // old str
	newStr    string // new replace str
}

// IsValidContent checks content type
func IsValidContent(ctype string) bool {
	for _, v := range CONTENT_TYPES {
		if strings.Contains(ctype, v) {
			return true
		}
	}
	return false
}

// New function return Proxy instance
func New(url, oldStr string, newStr string) *Proxy {
	return &Proxy{url, oldStr, newStr}
}

// proxy handler
func (p *Proxy) proxy(w http.ResponseWriter, r *http.Request) {
	// get url to proxy
	url := p.targetURL + r.URL.String()
	res, err := http.Get(url)
	if err != nil {
		log.Println("err fetching url", err)
		return
	}
	defer res.Body.Close()

	// read response data
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("err while parsing", err)
		return
	}

	if IsValidContent(res.Header.Get("Content-type")) {
		// Replace old str to new str
		data = bytes.Replace(data, []byte(p.oldStr), []byte(p.newStr), -1)
	}

	//write response
	w.Write(data)
}

func main() {
	// flags
	host := flag.String("host", "habrahabr.ru", "target host")
	oldStr := flag.String("old", "", "searching str")
	newStr := flag.String("new", "", "new replacement str")
	flag.Parse()
	//formatting host url
	*host = fmt.Sprintf("http://%v/", *host)

	// create proxy instance
	proxy := New(*host, *oldStr, *newStr)
	// proxy server address
	addr := "localhost:3000"
	// register url and handlerFunc called for "/" url
	http.HandleFunc("/", proxy.proxy)
	fmt.Println("Starting Server on", addr)
	// start server and check for errors
	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		// fatal and log error
		log.Fatal("Server could not started", err)
	}

}
