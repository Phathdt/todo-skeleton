package sdkcm

import (
	"errors"
	"net/http"
)

var (
	ErrRequestDataInvalid     = func(s string) *customError { return CustomError("ErrRequestDataInvalid", s) }
	ErrNoPermission           = CustomError("ErrNoPermission", "you don't have permission to access")
	ErrAccessTokenInvalid     = CustomError("ErrAccessTokenInvalid", "invalid access token")
	ErrAccessTokenInactivated = CustomError("ErrAccessTokenInactivated", "access token is disabled")
	ErrUserNotFound           = CustomError("ErrUserNotFound", "user not found or deactivated")
	ErrTopicNotFound          = CustomError("ErrTopicNotFound", "topic not found")
	ErrUserNameMinMaxLength   = CustomError("ErrUserNameMinMaxLength", "UserName must have length greater than 3 and less than 100")
	ErrPasswordMinMaxLength   = CustomError("ErrPasswordMinMaxLength", "Password must have length greater than 6 and less than 50")
)

var (
	// data not found sometime is not an error
	// but we need this type to decouple from db (errNotFound mongodb and gorm)
	ErrDataNotFound = errors.New("data not found")
)

var (
	ErrSignNotMatched = func() AppError {
		return NewAppErr(errors.New("sign not matched"), http.StatusBadRequest, "sign not matched").WithCode("sign_not_matched")
	}

	ErrClientNotMatched = func() AppError {
		return NewAppErr(errors.New("client not matched"), http.StatusBadRequest, "client not matched").WithCode("client_not_matched")
	}

	ErrCannotFetchData = func(err error) AppError {
		return NewAppErr(err, http.StatusBadRequest, "can not fetch data").WithCode("cannot_fetch_data")
	}
	ErrDB = func(err error) AppError {
		return NewAppErr(err, http.StatusBadRequest, "db error").WithCode("db_error")
	}
	ErrInvalidRequest = func(err error) AppError {
		return NewAppErr(err, http.StatusBadRequest, "invalid request").WithCode("invalid_request")
	}
	ErrInvalidRequestWithMessage = func(err error, message string) AppError {
		return NewAppErr(err, http.StatusBadRequest, message).WithCode("invalid_request")
	}
	ErrWithMessage = func(root error, err ErrorWithKey) AppError {
		if root == nil {
			return NewAppErr(errors.New(err.Error()), http.StatusBadRequest, err.Error()).WithCode(err.Key())
		}
		return NewAppErr(root, http.StatusBadRequest, err.Error()).WithCode(err.Key())
	}
	ErrCustom = func(root error, err ErrorWithKey) AppError {
		if root == nil {
			return NewAppErr(errors.New(err.Error()), http.StatusBadRequest, err.Error()).WithCode(err.Key())
		}
		return NewAppErr(root, http.StatusBadRequest, err.Error()).WithCode(err.Key())
	}
	ErrUnauthorized = func(root error, err ErrorWithKey) AppError {
		if root == nil {
			return NewAppErr(errors.New(err.Error()), http.StatusUnauthorized, err.Error()).WithCode(err.Key())
		}
		return NewAppErr(root, http.StatusUnauthorized, err.Error()).WithCode(err.Key())
	}
)

type ErrorWithKey interface {
	error
	Key() string
}

type AppError struct {
	// We don't show root cause to the clients
	RootCause  error  `json:"-"`
	Code       string `json:"code"`
	Log        string `json:"log"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewAppErr(err error, statusCode int, msg string) AppError {
	return AppError{RootCause: err, Log: err.Error(), StatusCode: statusCode, Message: msg}
}

// AppError is error
func (ae AppError) Error() string {
	return ae.Message
}

func (ae AppError) RootError() error {
	if root, ok := ae.RootCause.(AppError); ok {
		return root.RootError()
	}

	return ae.RootCause
}

func (ae AppError) WithCode(code string) AppError {
	ae.Code = code
	return ae
}

type customError struct {
	k string
	v string
}

func (ce *customError) Error() string {
	return ce.v
}

func (ce *customError) Key() string {
	return ce.k
}

func CustomError(k, v string) *customError {
	return &customError{k, v}
}
