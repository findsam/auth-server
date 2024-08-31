package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hashOne, errOne := HashPassword(password)
	hashTwo, errTwo := HashPassword(password)

	if errOne != nil || errTwo != nil {
		t.Errorf("error hashing a password")
	}

	if hashOne == password {
		t.Errorf("hash should not be equal to the plain password")
	}

	if hashOne == "" {
		t.Errorf("hash should not be empty")
	}

	if hashOne == hashTwo {
		t.Errorf("hashes should be different for the same input due to salting")
	}
}

func TestComparePasswords(t *testing.T) {
	password := "password"
	hashedPassword, _ := HashPassword(password)

	if !ComparePasswords(hashedPassword, []byte(password)) {
		t.Error("ComparePasswords should return true for matching password")
	}

	if ComparePasswords(hashedPassword, []byte("wrongPassword")) {
		t.Error("ComparePasswords should return false for non-matching password")
	}
}
