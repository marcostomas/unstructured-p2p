package server

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "helloer\n")
}

func SearchFlooding(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func SearchRandomWalk(w http.ResponseWriter, req *http.Request) {

}

func SearchInDepth(w http.ResponseWriter, req *http.Request) {
	
}

func InitServer() {

	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/SearchFlooding", SearchFlooding)
	http.HandleFunc("/SearchRandomWalk", SearchRandomWalk)
	http.HandleFunc("/SearchInDepth", SearchInDepth)

	http.ListenAndServe(":8090", nil)
}
