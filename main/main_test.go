package main_test

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"golang.org/x/net/http2"
)

func BenchmarkM(b *testing.B) {
	client := http.Client{
		Transport: &http2.Transport{
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
			AllowHTTP: true,
		},
	}

	res, err := client.Get("http://127.0.0.1:8888/foobar")
	check(err)
	defer res.Body.Close()
	if res.ProtoMajor != 2 {
		fmt.Println("proto not h2c")
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	check(err)
	b.Log(string(body))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
