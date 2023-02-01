package main

import (
	"net/http"
)

func HelloHandler0001(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func main0001() {
	// Create a request limiter per handler.
	http.Handle("/", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(1, nil), HelloHandler))
	http.ListenAndServe(":12345", nil)
}
