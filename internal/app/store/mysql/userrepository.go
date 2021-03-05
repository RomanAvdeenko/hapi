package mysql

import (
	"database/sql"
	"hapi/internal/app/model"
	"hapi/internal/app/store"
)

const (
	queryUserByID    = "SELECT user_id, user_login, user_password, user_decrypt_password, user_name, user_enabled, user_balance, user_group_id FROM user WHERE user_id = ?"
	queryUserByLogin = "SELECT user_id, user_login, user_password, user_decrypt_password, user_name, user_enabled, user_balance, user_group_id FROM user WHERE user_login = ?"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	u := new(model.User)

	err := r.store.db.QueryRow(queryUserByID, id).Scan(
		&u.ID, &u.Login, &u.EncryptedPassword, &u.DecryptedPassword, &u.Name, &u.Enabled, &u.Balance, &u.Privileges,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByLogin(login string) (*model.User, error) {
	u := new(model.User)
	err := r.store.db.QueryRow(queryUserByLogin, login).Scan(
		&u.ID, &u.Login, &u.EncryptedPassword, &u.DecryptedPassword, &u.Name, &u.Enabled, &u.Balance, &u.Privileges,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}
