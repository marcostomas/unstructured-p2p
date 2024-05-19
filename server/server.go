package server

import (
	"fmt"
	"net/http"
)

func HelloE(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "helloer\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func InitServer(endercco string, porta string) {

	http.HandleFunc("/hello", HelloE)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(endercco+":"+porta, nil)
}
