package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//Proxy struct
type Proxy struct {
	targetURL string // target host url
	oldStr    string // old str
	newStr    string // new replace str
}

func New(url, oldStr string, newStr string) *Proxy {
	return &Proxy{url, oldStr, newStr}
}

// proxy handler
func (p *Proxy) proxy(w http.ResponseWriter, r *http.Request) {
	// get url to proxy
	url := p.targetURL + r.URL.String()
	log.Println("url to fetch", url)
	res, err := http.Get(url)
	if err != nil {
		log.Println("err fetching url", err)
		return
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("err while parsing", err)
		return
	}

	data = bytes.Replace(data, []byte(p.oldStr), []byte(p.newStr), -1)
	w.Write(data)
}

func main() {

	// flags
	host := flag.String("host", "habrahabr.ru", "target host")
	oldStr := flag.String("old", "", "oldStr value")
	newStr := flag.String("new", "", "newStr value")
	flag.Parse()

	*host = fmt.Sprintf("http://%v/", *host)

	proxy := New(*host, *oldStr, *newStr)
	// proxy server address
	addr := "localhost:3000"
	// register url and handlerFunc called for "/" url
	http.HandleFunc("/", proxy.proxy)
	// write that server starting in console
	fmt.Println("Starting Server on", addr)
	// start server and check for errors
	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		// fatal and log error
		log.Fatal("Server could not started", err)
	}

}
