package apperrors

import "errors"

var (
	ErrNoDataFound = errors.New("no data")

	ErrProductNotFound = errors.New("product is not found")

	ErrUnExpectedError = errors.New("something went wrong") //1

	ErrAuthHeaderMissing = errors.New("authorization header is missing")

	ErrInvalidAuthToken = errors.New("auth token is invalid")

	ErrUnauthorized = errors.New("unauthorized")

	ErrBadRequest = errors.New("bad request")

	ErrInvalidProductID = errors.New("product id is invalid")

	ErrUserAlreadyExists = errors.New("user already exists")

	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrInvalidRequest = errors.New("invalid request body")

	ErrUserDoesNotExist = errors.New("user does not exist")

	ErrInvalidInput = errors.New("invalid input please check your input")
)
