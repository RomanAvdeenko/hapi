package mysql

import (
	"hapi/internal/app/model"
	"hapi/internal/app/store"
)

const (
	querySelectUserAccountByUserID = `
	SELECT 
		useraccount_id,
		useraccount_user_id,
		useraccount_ts,
		useraccount_date,
		useraccount_sum,
		useraccount_desc,
		useraccount_operator,
		useraccount_destination_id
	FROM useraccount
	WHERE
		useraccount_user_id = ?`
)

type UserAccountRepository struct {
	store *Store
}

func (r *UserAccountRepository) FindByUserID(userid uint) ([]model.UserAccount, error) {
	got := make([]model.UserAccount, 0, 4)

	rows, err := r.store.db.Query(querySelectUserAccountByUserID, userid)
	if err != nil {
		//fmt.Printf("Error Query: %v", err)
		return nil, store.ErrRecordNotFound
	}
	defer rows.Close()

	for rows.Next() {
		r := model.UserAccount{}
		err := rows.Scan(&r.ID, &r.UserID, &r.TimeStamp, &r.Date, &r.Sum, &r.Description, &r.Operator, &r.DestinationID)
		if err != nil {
			//fmt.Printf("Error Scan: %v", err)
			return nil, store.ErrRecordNotFound
		}
		got = append(got, r)
	}

	return got, nil
}
