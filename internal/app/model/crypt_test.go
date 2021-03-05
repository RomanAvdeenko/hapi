package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encryptString(t *testing.T) {
	u := TestUser(t)

	assert.Equal(t, encryptString(u.DecryptedPassword), u.EncryptedPassword)
	assert.NotEqual(t, encryptString(u.DecryptedPassword+u.DecryptedPassword), u.EncryptedPassword)
}
