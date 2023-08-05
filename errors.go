package cqrs

import "errors"

var (
	ErrHandlerMustBeNonNilFunction                          = errors.New("handler must be non nil function")
	ErrHandlerMustHaveExactTwoArguments                     = errors.New("handler must have exact 2 arguments")
	ErrFirstArgumentOfHandlerMustBeContext                  = errors.New("first argument of handler must be context.Context")
	ErrSecondArgumentOfHandlerMustBeStructOrPointerOfStruct = errors.New("second argument of handler must be struct or pointer of struct")
	ErrHandlerMustHaveExactOneResult                        = errors.New("handler must return exact 1 result")
	ErrHandlerResultMustBeError                             = errors.New("hander result must be error")
	ErrCommandAlreadyRegistered                             = errors.New("command already registered")
	ErrCommandHasNotRegisteredYet                           = errors.New("command has not registered yet")
)

var (
	ErrHandlerMustHaveExactTwoResults   = errors.New("handler must return exact 2 results")
	ErrSecondResultOfHandlerMustBeError = errors.New("second result of handler must be error")
	ErrQueryAlreadyRegistered           = errors.New("query already registered")
	ErrQueryHasNotRegisteredYet         = errors.New("query has not registered yet")
)

var (
	ErrEventProviderNotFound = errors.New("event provider not found")
)
