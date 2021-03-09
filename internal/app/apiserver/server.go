package apiserver

import (
	"encoding/json"
	"hapi/internal/app/model"
	"hapi/internal/app/store"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	logFile     = "./log"
	sessionName = "hitline_billing_session"
	sessionKey  = "user_id"

	ctxKeyUser ctxKey = iota
	ctxKeyUserID
	ctxKeyUUID
)

type ctxKey uint8

type server struct {
	router        *mux.Router
	logger        *logrus.Logger
	sessionsStore sessions.Store
	store         store.Store
}

// NewServer is constructor helper
func NewServer(config *Config, store store.Store, sessionsStore sessions.Store) *server {
	s := &server{
		router:        mux.NewRouter(),
		logger:        logrus.New(),
		store:         store,
		sessionsStore: sessionsStore,
	}

	if err := s.configureLogger(config); err != nil {

	}
	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	// Debug not found handler
	s.router.NotFoundHandler = http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugf("Not found: %v %v\n", r.Method, r.RequestURI)
		w.WriteHeader(http.StatusNotFound)
	}))

	s.router.Use(s.setRequestID)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.Use(s.logMiddleware)
	s.router.Use(s.panicMiddleware)

	s.router.HandleFunc("/", s.handleRoot())
	s.router.HandleFunc("/status", s.handleStatus())
	s.router.HandleFunc("/login", s.handleSessionCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/user", s.handleUser()).Methods(http.MethodGet)

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser, s.checkAccess)
	private.HandleFunc("/whoami", s.whoAmI()).Methods(http.MethodGet)
	private.HandleFunc("/userid/{userID:"+model.IDRegexp+"}/", s.handleUserByID()).Methods(http.MethodGet)
	private.HandleFunc("/user/{login:"+model.LoginRegexp+"}/", s.handleUserByLogin()).Methods(http.MethodGet)
	private.HandleFunc("/userid/{userID:"+model.IDRegexp+"}/account/", s.handleUserAccountByUserID()).Methods(http.MethodGet)
	private.HandleFunc("/userid/{userID:"+model.IDRegexp+"}/package/", s.handleUserPackageByUserID()).Methods(http.MethodGet)

}

func (s *server) configureLogger(config *Config) error {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		s.logger.Error("logger parse debug level")
		return err
	}

	s.logger.SetLevel(level)

	// //mw := io.MultiWriter(os.Stdout, logFile)
	// f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	// if err != nil {
	// 	s.logger.Error("open log file")
	// 	return err
	// }
	// s.logger.SetOutput(f)

	return nil
}

// error is a wraper for json erros respond
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// respond is a common server json respond
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.WriteHeader(code)
	if payload != nil {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "\t")
		enc.Encode(payload)
	}
}
