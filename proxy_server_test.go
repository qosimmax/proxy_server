package main

import (
	"log"
	"net/http"
	"testing"
)

func TestContenType(t *testing.T) {

	content_type := "text/html"
	if !IsValidContent(content_type) {
		t.Errorf("Invalid Content-Type:%v", content_type)
	}

}

func ProxyTest(t *testing.T) {
	res, err := http.Get("http://127.0.0.1:3000")
	if err != nil {
		t.Error(err)
	}
	log.Println(res)
}
