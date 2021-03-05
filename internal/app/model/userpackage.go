package model

import (
	"net"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type UserPackage struct {
	ID          uint   `json:"id"`
	UniqueID    uint   `json:"unique_id"`
	ParentID    uint   `json:"parent_id"`
	UserID      uint   `json:"-"`
	PackageID   uint   `json:"package_id"`
	StartDate   string `json:"start_date"`
	StopDate    string `json:"stop_date"`
	Enabled     bool   `json:"enabled"`
	ModifyDate  string `json:"modify_date"`
	IPAddress   net.IP `json:"ip_address,omitempty"`
	IPNetwork   uint   `json:"-"`
	MACAddress  string `json:"mac_address,omitempty"`
	Description string `json:"note,omitempty"`
}

func (u *UserPackage) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.ID, validation.Match(IDRE)),
		validation.Field(&u.UniqueID, validation.Match(IDRE)),
		validation.Field(&u.PackageID, validation.Match(IDRE)),
		validation.Field(&u.UserID, validation.Match(IDRE)),
		validation.Field(&u.PackageID, validation.Match(IDRE)),
		validation.Field(&u.StartDate, validation.Required),
		validation.Field(&u.StopDate, validation.Required),
		validation.Field(&u.Enabled, validation.Required),
		validation.Field(&u.ModifyDate, validation.Required),
		validation.Field(&u.IPAddress, is.IPv4),
		validation.Field(&u.IPNetwork, is.Int),
		validation.Field(&u.MACAddress, is.MAC),
	)
}
