package apiserver

import (
	"database/sql"
	"hapi/internal/app/store/mysql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

// Start ...
func Start(config *Config) error {
	db, err := NewDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := mysql.New(db)
	sessionsStore := sessions.NewCookieStore([]byte(config.SessionKey))
	sessionsStore.MaxAge(config.CookieAge)

	server := NewServer(config, store, sessionsStore)
	if err := server.configureLogger(config); err != nil {
		return err
	}

	serverAddress := config.BindAddr + ":" + config.Port
	server.logger.Infof("Starting server %s...", serverAddress)

	return http.ListenAndServe(serverAddress, server)
}

// NewDB create sql.DD for DI to Store
func NewDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
