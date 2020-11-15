package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func main() {
	// fastHttp Client
	req := &fasthttp.Request{}
	res := &fasthttp.Response{}

	req.SetRequestURI("http://0.0.0.0:32768/app/")
	req.Header.SetMethod("GET")

	err := fasthttp.Do(req, res)
	if err != nil {
		return
	}

	fmt.Println("")
}
