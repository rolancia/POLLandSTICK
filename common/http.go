package common

import (
	"context"
	"net/http"
	"sync"

	"github.com/rolancia/go-stick/stick"
)

type RequestCtxKey string

type Request struct {
	Req *http.Request
	Res http.ResponseWriter
}

func GetRequestFrom(ctx context.Context) Request {
	return stick.GetFrom[Request](ctx, RequestCtxKey(""))
}

func WithRequest(ctx context.Context, req Request) context.Context {
	return stick.With(ctx, RequestCtxKey(""), req)
}

//
type RequestWgCtxKey string

func GetRequestWgFrom(ctx context.Context) *sync.WaitGroup {
	return stick.GetFrom[*sync.WaitGroup](ctx, RequestWgCtxKey(""))
}

func WithRequestWg(ctx context.Context, wg *sync.WaitGroup) context.Context {
	return stick.With(ctx, RequestWgCtxKey(""), wg)
}
