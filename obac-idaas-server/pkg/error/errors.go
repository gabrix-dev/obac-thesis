package error

import (
	"reflect"
)

type CustomError struct {
	Error            int    `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type BadRequestError struct {
	err         error
	description string
}

type NotFoundError struct {
	err         error
	description string
}

type NotAuthorizedError struct {
	err         error
	description string
}

func (e BadRequestError) Error() string {
	return e.err.Error()
}

func (e NotFoundError) Error() string {
	return e.err.Error()
}

func (e NotAuthorizedError) Error() string {
	return e.err.Error()
}

func (e BadRequestError) GetDescription() string {
	return e.description
}

func (e NotFoundError) GetDescription() string {
	return e.description
}

func (e NotAuthorizedError) GetDescription() string {
	return e.description
}

func (e BadRequestError) Is(target error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(target)
}

func (e NotFoundError) Is(target error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(target)
}

func (e NotAuthorizedError) Is(target error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(target)
}

func NewBadRequestError(err error, description string) BadRequestError {
	return BadRequestError{err: err, description: description}
}

func NewNotFoundError(err error, description string) NotFoundError {
	return NotFoundError{err: err, description: description}
}

func NewNotAuthorizedError(err error, description string) NotAuthorizedError {
	return NotAuthorizedError{err: err, description: description}
}
