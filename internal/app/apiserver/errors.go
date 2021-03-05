package apiserver

import "errors"

var (
	errIncorrectLoginOrPassword = errors.New("incorrect login or password")
	errNotAuthenticated         = errors.New("not authenticated")
	errValidation               = errors.New("incorrect login or password_")
	errBadRequest               = errors.New("bad request")
	errNotFound                 = errors.New("not found")
	errAccessDeny               = errors.New("access deny")
)
