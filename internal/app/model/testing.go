package model

import (
	"regexp"
	"testing"
)

var (
	LoginRE = regexp.MustCompile(LoginRegexp)
	IDRE    = regexp.MustCompile(IDRegexp)
)

// TestUser helper
func TestUser(t *testing.T) *User {
	return &User{
		ID:                3003,
		Login:             "0962964374",
		EncryptedPassword: "ba15c212eeaeb515f73e3e1a198bff7e54a9cbd1",
		DecryptedPassword: "ravvar5",
		Name:              "Иван Иванов",
		Enabled:           true,
		Balance:           459.25,
	}
}
