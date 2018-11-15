package session

import (
	"fmt"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	from := ConstFromC
	srcID := 2
	managerID := 0
	userID := 0
	duration := time.Duration(6) * time.Hour * 100
	secretKey := "d61de49e5303f036692b171fe0d279e8"

	getSecretKey := func(int, int) (string, error) { return secretKey, nil }

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
