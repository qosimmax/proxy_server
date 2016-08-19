package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

// Content types
var CONTENT_TYPES = []string{"html", "text", "json", "xml"}

// regexp Compile
var RP = regexp.MustCompile("[\\p{L}\\d_]+")

// IsValidContent function check the Content-Type header
func IsValidContent(ctype string) bool {
	for _, v := range CONTENT_TYPES {
		if strings.Contains(ctype, v) {
			return true
		}
	}
	return false
}

// replace old str to new str
func Replace(data *[]byte, oldStr []byte, newStr []byte) {
	// use a unicode character class to include the digits and underscore
	*data = RP.ReplaceAllFunc(*data, func(b []byte) []byte {
		if bytes.Equal(b, oldStr) {
			return newStr
		}
		return b
	})
}

// Decode gzip data
func DecodeZip(b []byte) (z []byte, err error) {
	reader, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return z, err
	}

	z, err = ioutil.ReadAll(reader)
	return
}

// gzip data
func EncodeZip(b []byte) (z []byte, err error) {
	var bz bytes.Buffer
	gz := gzip.NewWriter(&bz)
	if _, err = gz.Write(b); err != nil {
		return
	}
	if err = gz.Flush(); err != nil {
		return
	}
	if err = gz.Close(); err != nil {
		return
	}
	z = bz.Bytes()
	return
}

type transport struct {
	http.RoundTripper
	oldStr string
	newStr string
}

// RoundTripper
func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Close Body
	defer resp.Body.Close()

	// check the Content-Type header
	if !IsValidContent(resp.Header.Get("Content-type")) {
		return resp, nil
	}

	//read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil
	}

	// check the Content-Encoding
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		//decode gzip data
		body, err = DecodeZip(body)
		if err != nil {
			return resp, nil
		}
		// Replace old str with new
		Replace(&body, []byte(t.oldStr), []byte(t.newStr))

		// encode to gzip
		body, err = EncodeZip(body)
		if err != nil {
			return resp, nil
		}

	default:
		// Replace old str with new
		Replace(&body, []byte(t.oldStr), []byte(t.newStr))

	}

	// set body and ContentLength
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp.ContentLength = int64(len(body))

	return resp, nil
}

// Prox object
type Prox struct {
	// target url of reverse proxy
	target *url.URL
	// instance of Go ReverseProxy
	proxy *httputil.ReverseProxy
}

// New function create instance of Prox
func New(url *url.URL) *Prox {
	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	r.Host = r.URL.Host
	p.proxy.ServeHTTP(w, r)
}

func main() {
	// constants
	const (
		Port               = ":3000"
		defaultTarget      = "habrahabr.ru"
		defaultTargetUsage = "default redirect url, '127.0.0.1:8080'"
		defaultOldStr      = "Go"
		defaultOldStrUsage = "default value, 'Go'"
		defaultNewStr      = "Golang"
		defaultNewStrUsage = "default new replacer value, 'Golang'"
	)

	// flags
	host := flag.String("host", defaultTarget, defaultTargetUsage)
	oldStr := flag.String("old", defaultOldStr, defaultOldStrUsage)
	newStr := flag.String("new", defaultNewStr, defaultNewStrUsage)
	flag.Parse()

	//checking host url
	*host = fmt.Sprintf("https://%v/", *host)
	url, err := url.Parse(*host)
	if err != nil {
		log.Fatal("Error parsing URL")
	}

	fmt.Println("server listen on:", Port)
	fmt.Println("redirecting to:", *host)
	// proxy instance
	proxy := New(url)
	proxy.proxy.Transport = &transport{http.DefaultTransport, *oldStr, *newStr}
	http.HandleFunc("/", proxy.handle)
	// start proxy server
	http.ListenAndServe(Port, nil)
}
