package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hapi/internal/app/model"
	"hapi/internal/app/store/mysql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

const (
	envDatabaseURL = "DATABASE_URL"
)

var (
	databaseURL = "api:jwjii3ke8dnsjGdsx0@tcp(10.2.1.20:3306)/billing"
	hashKey     = []byte("secter_key")
)

func TestMain(m *testing.M) {
	if env := os.Getenv(envDatabaseURL); env != "" {
		databaseURL = env
	}
	os.Exit(m.Run())
}

func TestServer_AuthenticateUser(t *testing.T) {
	u := model.TestUser(t)
	//log.Println("testUser():", u)

	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expextedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expextedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			cookieValue:  nil,
			expextedCode: http.StatusUnauthorized,
		},
	}

	coockie := securecookie.New(hashKey, nil)
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	store := mysql.TestStore(t, databaseURL)
	server := NewServer(&Config{}, store, sessions.NewCookieStore(hashKey))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieString, _ := coockie.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieString))

			server.authenticateUser(fakeHandler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expextedCode, rec.Code)
		})
	}
}

func TestServer_handleSessionCreate(t *testing.T) {
	u := model.TestUser(t)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"login":    u.Login,
				"password": u.DecryptedPassword,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "OK",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "bad login or password",
			payload: map[string]string{
				"login":    u.Login,
				"password": strings.Join([]string{u.DecryptedPassword, u.DecryptedPassword}, "!"),
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := mysql.TestStore(t, databaseURL)
			s := NewServer(&Config{}, store, sessions.NewCookieStore([]byte(hashKey)))

			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodPost, "/login", b)
			rec := httptest.NewRecorder()

			s.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

		})
	}
}
