package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Content types
var CONTENT_TYPES = []string{"html", "text", "json", "xml"}

//Proxy struct
type Proxy struct {
	targetURL string // target host url
	oldStr    []byte // old str
	newStr    []byte // new replace str
}

// IsValidContent function check the Content-Type header
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
	return &Proxy{url, []byte(oldStr), []byte(newStr)}
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

	// check the Content-Type header
	if !IsValidContent(res.Header.Get("Content-type")) {
		// write response
		w.Write(data)
		return
	}

	// replace old str to new str
	// use a unicode character class to include the digits and underscore
	rp := regexp.MustCompile("[\\p{L}\\d_]+")
	data = rp.ReplaceAllFunc(data, func(b []byte) []byte {
		if bytes.Equal(b, p.oldStr) {
			return p.newStr
		}
		return b
	})

	// write response
	w.Write(data)
}

func main() {
	// flags
	host := flag.String("host", "habrahabr.ru", "target host")
	oldStr := flag.String("old", "", "searching str")
	newStr := flag.String("new", "", "new replacement str")

	flag.Parse()
	//checking host url
	*host = fmt.Sprintf("http://%v/", *host)
	_, err := url.Parse(*host)
	if err != nil {
		log.Fatal("Error parsing URL")
	}

	// create proxy instance
	proxy := New(*host, *oldStr, *newStr)
	// proxy server address
	addr := "localhost:3000"
	// register url and handlerFunc called for "/" url
	http.HandleFunc("/", proxy.proxy)
	fmt.Println("Starting Server on", addr)
	// start server and check for errors
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		// fatal and log error
		log.Fatal("Server could not started", err)
	}

}
