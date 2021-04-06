package apiserver

import (
	"encoding/json"
	"fmt"
	"hapi/internal/app/model"
	"io"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

func (s *server) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "")
	}
}
func (s *server) handleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}
}

func (s *server) handleSessionCreate() http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(request)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		// Validate
		if err := validation.ValidateStruct(req, validation.Field(&req.Login, validation.Match(model.LoginRE))); err != nil {
			s.error(w, r, http.StatusUnauthorized, errValidation)
			return
		}
		// Try to find user in DB
		u, err := s.store.User().FindByLogin(req.Login)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectLoginOrPassword)
			return
		}

		session, err := s.sessionsStore.New(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values[sessionKey] = u.ID

		if err := s.sessionsStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) whoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		s.respond(w, r, http.StatusOK, u.Login)
	}
}

func (s *server) handleUserByLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVars := mux.Vars(r)

		//fmt.Println("urlVArs:", urlVars)

		login, ok := urlVars["login"]
		//fmt.Println("login:", login)
		if !ok {
			s.error(w, r, http.StatusBadRequest, errBadRequest)
			return
		}
		u, err := s.store.User().FindByLogin(login)

		if err != nil {
			//fmt.Println("Err: ", err)
			s.respond(w, r, http.StatusNotFound, map[string]string{"error": errNotFound.Error()})
			return
		}
		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handleUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVars := mux.Vars(r)

		userID, ok := urlVars["userID"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, errBadRequest)
			return
		}

		id, err := strconv.Atoi(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		u, err := s.store.User().FindByID(uint(id))

		if err != nil {
			s.respond(w, r, http.StatusNotFound, map[string]string{"error": errNotFound.Error()})
			return
		}
		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handleUserAccountByUserID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVars := mux.Vars(r)

		userID, ok := urlVars["userID"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, errBadRequest)
			return
		}

		id, err := strconv.Atoi(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		res, err := s.store.UserAccount().FindByUserID(uint(id))

		if err != nil {
			s.respond(w, r, http.StatusOK, map[string]string{"error": errNotFound.Error()})
			return
		}
		s.respond(w, r, http.StatusOK, res)
	}
}

func (s *server) handleUserPackageByUserID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVars := mux.Vars(r)

		userID, ok := urlVars["userID"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, errBadRequest)
			return
		}

		id, err := strconv.Atoi(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		res, err := s.store.UserPackage().FindByUserID(uint(id))

		if err != nil {
			fmt.Println("Err handleUserPackageByUserID(): ", err)
			s.respond(w, r, http.StatusOK, map[string]string{"error": errNotFound.Error()})
			return
		}
		s.respond(w, r, http.StatusOK, res)
	}
}

// handleUser() gets user based on session information
func (s *server) handleUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		session, err := s.sessionsStore.Get(r, sessionName)
		if err != nil {
			s.error(rw, r, http.StatusInternalServerError, err)
			return
		}

		userID, ok := session.Values[sessionKey]

		if !ok {
			s.error(rw, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		user, err := s.store.User().FindByID(userID.(uint))

		if err != nil {
			s.error(rw, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		s.respond(rw, r, http.StatusOK, user)
	}
}
