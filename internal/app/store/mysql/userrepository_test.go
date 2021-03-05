package mysql_test

import (
	"testing"

	"hapi/internal/app/model"
	"hapi/internal/app/store"
	"hapi/internal/app/store/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Find(t *testing.T) {
	userID := uint(0)
	s := mysql.TestStore(t, databaseURL)
	_, err := s.User().FindByID(userID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u1, err := s.User().FindByID(u.ID)
	//fmt.Printf("user: %+v\n", u1)
	assert.NoError(t, err)
	assert.Equal(t, u1.ID, u.ID)
	assert.Equal(t, u1.Login, u.Login)

	userLogin := "NonExixtingLogin_734hhdyualjclsjjuu7hdkdbkky6dfh"
	_, err = s.User().FindByLogin(userLogin)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
}
