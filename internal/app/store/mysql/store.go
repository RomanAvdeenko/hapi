package mysql

import (
	"database/sql"
	"hapi/internal/app/store"
	"time"
)

const queryTimeout = 1 * time.Second

type Store struct {
	//ctx            context.Context
	db                    *sql.DB
	userRepository        *UserRepository
	userAccountRepository *UserAccountRepository
	userPackageRepository *UserPackageRepository
}

func New(db *sql.DB) *Store {
	s := &Store{db: db}
	u := &UserRepository{store: s}

	s.userRepository = u

	uA := &UserAccountRepository{store: s}
	s.userAccountRepository = uA

	uP := &UserPackageRepository{store: s}
	s.userPackageRepository = uP

	return s
}

func (s *Store) User() store.UserRepository {
	return s.userRepository
}

func (s *Store) UserAccount() store.UserAccountRepository {
	return s.userAccountRepository
}

func (s *Store) UserPackage() store.UserPackageRepository {
	return s.userPackageRepository
}
