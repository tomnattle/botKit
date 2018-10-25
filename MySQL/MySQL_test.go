package MySQL

import (
	"testing"
)

func TestMySQL(t *testing.T) {
	conn, err := GetConn()
	if err != nil {
		t.Fatal(err)
	}

	if err := conn.Ping(); err != nil {
		t.Fatal(err)
	}

	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
}
