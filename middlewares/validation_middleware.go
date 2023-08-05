package middlewares

import (
	"context"

	"github.com/vulpes-ferrilata/cqrs"
	"github.com/vulpes-ferrilata/cqrs/pkg/validator"
)

func NewValidationMiddleware(validate validator.Validate) *ValidationMiddleware {
	return &ValidationMiddleware{
		validate: validate,
	}
}

type ValidationMiddleware struct {
	validate validator.Validate
}

func (v ValidationMiddleware) CommandMiddleware() cqrs.CommandMiddlewareFunc {
	return func(handler cqrs.CommandHandlerFunc[any]) cqrs.CommandHandlerFunc[any] {
		return func(ctx context.Context, command any) error {
			if err := v.validate.StructCtx(ctx, command); err != nil {
				return err
			}

			if err := handler(ctx, command); err != nil {
				return err
			}

			return nil
		}
	}
}

func (v ValidationMiddleware) QueryMiddleware() cqrs.QueryMiddlewareFunc {
	return func(handler cqrs.QueryHandlerFunc[any, any]) cqrs.QueryHandlerFunc[any, any] {
		return func(ctx context.Context, command any) (interface{}, error) {
			if err := v.validate.StructCtx(ctx, command); err != nil {
				return nil, err
			}

			result, err := handler(ctx, command)
			if err != nil {
				return nil, err
			}

			return result, nil
		}
	}
}
