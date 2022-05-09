package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cloudwego/netpoll"
	http2 "github.com/cloudwego/netpoll-http2"
	"github.com/rolancia/go-stick/stick"

	"github.com/rolancia/POLLandSTICK/common"
	"github.com/rolancia/POLLandSTICK/sticks"
)

type chReqCtxKey string

func main() {
	appContext := context.Background()

	listener, err := netpoll.CreateListener("tcp", "127.0.0.1:8888")
	check(err)

	server := http2.Server{Handler: &serverHandler{}, IdleTimeout: time.Minute}
	chReq := make(chan context.Context, 1000)
	eventLoop, err := netpoll.NewEventLoop(server.ServeConn,
		netpoll.WithIdleTimeout(10*time.Minute),
		netpoll.WithOnPrepare(func(conn netpoll.Connection) context.Context {
			_ = conn.SetReadTimeout(3 * time.Minute)
			ctx := context.Background()
			ctx = stick.With(ctx, chReqCtxKey(""), chReq)
			return ctx
		}),
	)
	check(err)

	go func() {
		err = eventLoop.Serve(listener)
		check(err)
	}()

	httpBunch := stick.Bunch().
		Defer(func(ctx context.Context, err error) context.Context {
			if v := recover(); v != nil {
				fmt.Println(v)
			}
			common.GetRequestWgFrom(ctx).Done()
			return ctx
		}).
		I(sticks.HttpCheckerStick{}).
		L(sticks.BasicLogBunch).
		I(sticks.EchoStick{})

	appContext = stick.WithConfig(appContext, stick.Config{
		Worker: func(job func()) {
			go job()
		},
	})
	err = stick.Spin(appContext, chReq, httpBunch)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type serverHandler struct{}

func (sh *serverHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ch := stick.GetFrom[chan context.Context](req.Context(), chReqCtxKey(""))
	r := common.Request{Req: req, Res: w}
	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx := req.Context()
	ctx = common.WithRequest(ctx, r)
	ctx = common.WithRequestWg(ctx, &wg)
	ch <- ctx
	wg.Wait()
}
