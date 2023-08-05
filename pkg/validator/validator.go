package validator

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock_$GOPACKAGE
import "context"

type Validate interface {
	StructCtx(ctx context.Context, s interface{}) error
}
