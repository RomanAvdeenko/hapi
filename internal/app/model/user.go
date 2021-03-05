package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// User ...
type User struct {
	ID                uint    `json:"id"`
	Login             string  `json:"login"`
	EncryptedPassword string  `json:"password,omitempty"`
	Name              string  `json:"name"`
	Enabled           bool    `json:"enabled"`
	Balance           float64 `json:"balance"`
	DecryptedPassword string  `json:"-"`
	Privileges        uint    `json:"-"`
}

// Validate ...
func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.ID, validation.Required),
		validation.Field(&u.Login, validation.Match(LoginRE)),
	)
}
