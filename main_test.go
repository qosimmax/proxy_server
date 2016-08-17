package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// create proxy instance
	host := "https://habrahabr.ru"
	oldStr := "google"
	newStr := "Google"
	proxy := New(host, oldStr, newStr)
	// proxy server listen address
	addr := ":3000"
	// register url and handlerFunc called for "/" url
	http.HandleFunc("/", proxy.proxy)
	fmt.Println("Starting Server on", addr)
	// start server and check for errors
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		// fatal and log error
		log.Fatal("Server could not started", err)
	}
	os.Exit(m.Run())
}
