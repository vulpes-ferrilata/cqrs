package eventbus

import "errors"

var (
	ErrMiddlewareFuncMustNotBeNil                   = errors.New("middlewareFunc must not be nil")
	ErrHandlerFuncMustBeNonNilFunction              = errors.New("handlerFunc must be non nil function")
	ErrHandlerFuncMustHaveExactTwoArguments         = errors.New("handlerFunc must have exact 2 arguments")
	ErrFirstArgumentOfHandlerMustBeContext          = errors.New("first argument of handler must be context.Context")
	ErrSecondArgumentOfHandlerMustBePointerOfStruct = errors.New("second argument of handler must be pointer of struct")
	ErrHandlerFuncMustHaveExactOneResult            = errors.New("handlerFunc must return exact 1 result")
	ErrResultMustBeError                            = errors.New("result must be error")
	ErrEventMustBeNonNilPointerOfStruct             = errors.New("event must be non nil pointer of struct")
)