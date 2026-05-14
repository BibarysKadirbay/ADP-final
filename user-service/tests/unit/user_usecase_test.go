package unit

import (
	"testing"
)

func TestUserName(t *testing.T) {

	name := "Bibarys"

	if name == "" {
		t.Fatal("user name is empty")
	}
}

func TestUserEmail(t *testing.T) {

	email := "bibarys@test.com"

	if email == "" {
		t.Fatal("user email is empty")
	}
}
