package errorrs

import (
	"net/http"
)

var errstore map[int]string

func init() {
	errstore = map[int]string{
		CodeTokenExpired:               "access token expired",
		CodeTokenCreationFailed:        "creating access token failed",
		CodeDatabaseQueryBuilderFailed: "something went wrong",
		CodeDatabaseInsertFailed:       "failed to issue token",
		CodeAccessBlocked:              "you no longer have access. please login",
		CodeRefreshTokenExpired:        "you are not allowed to access this resource. please login",
		CodeUpdateQueryFailed:          "something went wrong",
	}
}

type ApplicationError struct {
	err  error
	code int
}

func InternalError(err error, code int) *ApplicationError {
	return &ApplicationError{err, code}
}

func (ApplicationError) Status() int {
	return http.StatusInternalServerError
}

func (a *ApplicationError) Error() string {
	return errstore[a.code]
}

func (a *ApplicationError) Trace() error {
	return a.err
}

type NotFoundError struct {
	*ApplicationError
}

func (NotFoundError) Status() int {
	return http.StatusNotFound
}

func NotFound(err error, code int) NotFoundError {
	return NotFoundError{InternalError(err, code)}
}

type UnAuthorizedError struct {
	*ApplicationError
}

func (UnAuthorizedError) Status() int {
	return http.StatusUnauthorized
}

func UnAuthorized(err error, code int) UnAuthorizedError {
	return UnAuthorizedError{InternalError(err, code)}
}
