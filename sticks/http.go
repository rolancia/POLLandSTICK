package sticks

import (
	"context"

	"github.com/rolancia/go-stick/stick"

	"github.com/rolancia/POLLandSTICK/common"
)

var _ stick.Stick = &HttpCheckerStick{}

type HttpCheckerStick struct {
}

func (h HttpCheckerStick) Ignore(_ context.Context, _ error) bool {
	return false
}

func (h HttpCheckerStick) Handle(ctx context.Context, err error) (context.Context, error) {
	if !stick.Has(ctx, common.RequestCtxKey("")) {
		panic("not http")
	}
	return ctx, err
}

var _ stick.Stick = &EchoStick{}

type EchoStick struct {
}

func (e EchoStick) Ignore(ctx context.Context, err error) bool {
	return err != nil
}

func (e EchoStick) Handle(ctx context.Context, _ error) (context.Context, error) {
	req := common.GetRequestFrom(ctx)
	path := req.Req.URL.Path
	_, err := req.Res.Write([]byte(path))
	return ctx, err
}
