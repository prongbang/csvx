package csvx_test

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/prongbang/csvx"
)

func TestFileHeaderReader(t *testing.T) {
	// Given
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	file, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write([]byte("ID,Name\n1,N1\n2,N2\n"))
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest("POST", "/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if err := request.ParseMultipartForm(1024); err != nil {
		t.Fatal(err)
	}
	fileHeader := request.MultipartForm.File["file"][0]

	// When
	result, err := csvx.FileHeaderReader(fileHeader)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 || result[0][0] != "1" || result[1][1] != "N2" {
		t.Error("FileHeaderReader error:", result)
	}
}

func TestFileHeaderReaderInvalidFile(t *testing.T) {
	_, err := csvx.FileHeaderReader(&multipart.FileHeader{})
	if err == nil {
		t.Error("FileHeaderReader should return error")
	}
}

func TestFileHeaderReaderEmptyFile(t *testing.T) {
	// Given
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_, err := writer.CreateFormFile("file", "empty.csv")
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest("POST", "/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if err := request.ParseMultipartForm(1024); err != nil {
		t.Fatal(err)
	}
	fileHeader := request.MultipartForm.File["file"][0]

	// When
	result, err := csvx.FileHeaderReader(fileHeader)

	// Then
	if err == nil {
		t.Error("FileHeaderReader empty file should return error")
	}
	if len(result) != 0 {
		t.Error("FileHeaderReader empty file rows error:", result)
	}
}
