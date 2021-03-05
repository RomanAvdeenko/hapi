package store

import "hapi/internal/app/model"

// UserRepository...
type UserRepository interface {
	//Create(u *model.User) error
	//FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByLogin(login string) (*model.User, error)
}

type UserAccountRepository interface {
	FindByUserID(uint) ([]model.UserAccount, error)
}

type UserPackageRepository interface {
	FindByUserID(uint) ([]model.UserPackage, error)
}
