package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type UserAccount struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"-"`
	TimeStamp     string  `json:"ts"`
	Date          string  `json:"-"`
	Sum           float64 `json:"sum"`
	Description   string  `json:"note"`
	Operator      uint    `json:"operator"`
	DestinationID uint    `json:"dest_id"`
}

func (u *UserAccount) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.ID, validation.Match(IDRE)),
		validation.Field(&u.UserID, validation.Match(IDRE)),
		validation.Field(&u.TimeStamp, validation.Required),
		validation.Field(&u.Date, validation.Required),
	)
}
