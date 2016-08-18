package main

import (
	"bytes"
	"testing"
)

func TestContenType(t *testing.T) {

	content_type := "text/html"
	if !IsValidContent(content_type) {
		t.Errorf("Invalid Content-Type:%v", content_type)
	}

}

func ReplaceTest(t *testing.T) {
	data := []byte("Hello world")
	oldStr := []byte("world")
	newStr := []byte("World")

	olddata := data
	Replace(&data, oldStr, newStr)
	if bytes.Equal(olddata, data) {
		t.Errorf("Error on Replace old:%v, new:%v", oldStr, newStr)
	}

}

func DecodeZipTest(t *testing.T) {
	data := []byte("Hello world")
	_, err := DecodeZip(data)
	if err != nil {
		t.Error("Decode gzip data error")
	}

}

func EncodeZipTest(t *testing.T) {
	data := []byte("Hello world")
	_, err := EncodeZip(data)
	if err != nil {
		t.Error("Encode data error")
	}

}
