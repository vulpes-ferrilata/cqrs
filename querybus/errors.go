package querybus

import "errors"

var (
	ErrMiddlewareFuncMustNotBeNil                   = errors.New("middlewareFunc must not be nil")
	ErrHandlerFuncMustBeNonNilFunction              = errors.New("handlerFunc must be non nil function")
	ErrHandlerFuncMustHaveExactTwoArguments         = errors.New("handlerFunc must have exact 2 arguments")
	ErrFirstArgumentOfHandlerMustBeContext          = errors.New("first argument of handler must be context.Context")
	ErrSecondArgumentOfHandlerMustBePointerOfStruct = errors.New("second argument of handler must be pointer of struct")
	ErrHandlerFuncMustHaveExactTwoResults           = errors.New("handlerFunc must return exact 2 results")
	ErrFirstResultOfHandlerMustBePointerOfStruct    = errors.New("first result of handler must be pointer of struct")
	ErrSecondResultOfHandlerMustBeError             = errors.New("second result of handler must be error")
	ErrQueryIsAlreadyRegistered                     = errors.New("query is already registered")
	ErrQueryMustBeNonNilPointerOfStruct             = errors.New("query must be non nil pointer of struct")
	ErrQueryHasNotRegisteredYet                     = errors.New("query has not registered yet")
)
