package session

import (
	"fmt"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	from := ConstFromA
	srcID := 1
	managerID := 1
	userID := 1
	duration := time.Duration(6) * time.Hour
	getSecretKey := func(int, int) (string, error) {
		secretKey := "abc"
		return secretKey, nil
	}

	s, err := GenerateSession(from, srcID, managerID, userID, duration, getSecretKey)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("session string is ", s)
	sBasic, err := VerifySession(from, s, getSecretKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("session struct is ", sBasic)
}
