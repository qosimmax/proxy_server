package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestContenType(t *testing.T) {

	content_type := "text/html"

	if !IsValidContent(content_type) {
		t.Errorf("Invalid Content-Type:%v", content_type)
	}

}

func TestReplace(t *testing.T) {
	data := []byte("Hello world")
	oldStr := []byte("world")
	newStr := []byte("World")

	olddata := data
	Replace(&data, oldStr, newStr)
	if bytes.Equal(olddata, data) {
		t.Errorf("Error on Replace old:%v, new:%v", oldStr, newStr)
	}

}

func TestDecodeZip(t *testing.T) {

	data := []byte("Hello world")

	gzdata, err := EncodeZip(data)
	if err != nil {
		t.Errorf("Encode data error:%v", err)
		return
	}

	_, err = DecodeZip(gzdata)
	if err != nil {
		t.Errorf("Decode gzip data error:%v", err)
	}

}

func EncodeZipTest(t *testing.T) {
	data := []byte("Hello world")

	_, err := EncodeZip(data)
	if err != nil {
		t.Error("Encode data error")
	}

}

func TestHandler(t *testing.T) {
	oldStr := "Go"
	newStr := "Golang"

	tr := &transport{http.DefaultTransport, oldStr, newStr}
	c := &http.Client{Transport: tr}
	resp, err := c.Get("http://habrahabr.ru")
	if err != nil {
		t.Errorf("http request error:%v", err)
		return
	}
	//read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("read response error:%v", err)
		return
	}
	if !bytes.Contains(body, []byte(newStr)) {
		t.Errorf("'%v' not replaced", newStr)
	}
}
