package sticks

import (
	"context"
	"fmt"

	"github.com/rolancia/go-stick/stick"

	"github.com/rolancia/POLLandSTICK/common"
)

var BasicLogBunch = stick.Bunch().
	Defer(func(ctx context.Context, err error) context.Context {
		if !errorStick.Ignore(ctx, err) {
			errorStick.Handle(ctx, err)
		}
		return ctx
	}).
	I(AccessLogStick{})

//
var errorStick stick.Stick = &ErrorStick{}

type ErrorStick struct {
}

func (e ErrorStick) Ignore(_ context.Context, err error) bool {
	return err == nil
}

func (e ErrorStick) Handle(ctx context.Context, err error) (context.Context, error) {
	fmt.Println("Error:", err)
	return ctx, err
}

//
var _ stick.Stick = &AccessLogStick{}

type AccessLogStick struct {
}

func (a AccessLogStick) Ignore(ctx context.Context, err error) bool {
	return err != nil
}

func (a AccessLogStick) Handle(ctx context.Context, _ error) (context.Context, error) {
	req := common.GetRequestFrom(ctx)
	fmt.Println(req.Req.URL.Path)
	return ctx, nil
}
