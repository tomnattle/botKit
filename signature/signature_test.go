package signature

import (
	"bytes"
	"net/http"
	"testing"
)

func TestAddSignature(t *testing.T) {
	req, err := http.NewRequest("POST", "127.0.0.1", &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	err = AddSignature(req)
	if err != nil {
		t.Fatal(err)
	}
	pass, err := VerifySignature(req)
	if err != nil {
		t.Fatal(err)
	}
	if !pass {
		t.Fatal("not equal")
	}
}
