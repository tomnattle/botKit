package util

import (
	"testing"
)

func TestCheckPhone(t *testing.T) {
	success := []string{
		"12345678901",
		"16734532454",
		"12457563467",
	}

	for _, phone := range success {
		err := CheckPhone(phone)
		if err == nil {
			continue
		}
		t.Fatalf("CheckPhone error %s %v",
			phone, err)
	}

	fail := []string{
		"",
		"1",
		"a",
		"1245753467",
		"124575a3467",
		"124575323467",
		"a12345678901",
		"12345678901a",
		"a12345678901a",
	}

	for _, phone := range fail {
		err := CheckPhone(phone)
		if err != nil {
			continue
		}
		t.Fatalf("CheckPhone error %s pass",
			phone)
	}
}
