package SMS

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestSMS(t *testing.T) {
	err := Request("17602110192",
		fmt.Sprintf(cfg.Model, "123", 1))
	if err != nil {
		t.Fatal(err)
	}
}

func TestStepSMS(t *testing.T) {
	phone := "17602110192"
	content := fmt.Sprintf(cfg.Model, "123", 1)
	req, err := payload(phone, content)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
