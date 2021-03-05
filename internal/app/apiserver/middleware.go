package apiserver

import (
	"context"
	"fmt"
	"hapi/internal/app/model"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *server) panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err, ok := recover().(error); ok {
					s.error(w, r, http.StatusInternalServerError, fmt.Errorf("panic recover: %w", err))
				}
			}()
			next.ServeHTTP(w, r)
		})
}

func (s *server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &ResponseWriter{w, http.StatusOK}

			next.ServeHTTP(rw, r)

			s.logger.Debugf("%v %v %v <%v> (%s)\n", r.Context().Value(ctxKeyUUID), r.Method, r.RequestURI, rw.statusCode, time.Since(start))
		})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := s.sessionsStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		userID, ok := session.Values[sessionKey]

		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		user, err := s.store.User().FindByID(userID.(uint))

		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
		//next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUserID, userID)))
	})
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			uuid := uuid.New().String()
			w.Header().Set("X-Request-ID", uuid)
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUUID, uuid)))
		})
}

// Access control
func (s *server) checkAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			u := r.Context().Value(ctxKeyUser).(*model.User)
			if u.Privileges != 2 {
				s.error(w, r, http.StatusForbidden, errAccessDeny)
				return
			}
			next.ServeHTTP(w, r)
		})
}

// It's not work. May be gorolla lost main context?
// func (s *server) withTimeoutMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Nanosecond)
// 			defer cancel()

// 			r = r.WithContext(ctx)
// 			s.logger.Debug("withTimeoutMiddleware IN")
// 			next.ServeHTTP(w, r)
// 			s.logger.Debug("withTimeoutMiddleware OUT")
// 		})
// }
// //http.TimeoutHandler(next, 10*time.Nanosecond, "timeout")
