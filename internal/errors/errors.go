package errors

import "errors"

var (
	ErrInvalidHtppMethod                = errors.New("invalid http method")
	ErrInvalidJson                      = errors.New("invalid JSON")
	ErrInternalServerError              = errors.New("internal server error")
	ErrAccessDenied                     = errors.New("access denied")
	ErrInvalidPvzIdFormat               = errors.New("invalid pvz id format")
	ErrPvzDoesNotExist                  = errors.New("pvz does not exist")
	ErrReceptionInProgressDoesNotExist  = errors.New("reception in progress does not exist")
	ErrNoProductToDelete                = errors.New("no product to delete")
	ErrReceptionInProgressAlreadyExists = errors.New("reception in progress already exist")
)
