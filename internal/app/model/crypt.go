package model

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
)

// encryptString returns 3 times sha1()
func encryptString(s string) string {
	for i := 0; i < 3; i++ {
		h := sha1.New()
		io.WriteString(h, s)
		s = hex.EncodeToString(h.Sum(nil))
	}
	return s
}

func (u *User) ComparePassword(pass string) bool {
	return encryptString(pass) == u.EncryptedPassword
}
