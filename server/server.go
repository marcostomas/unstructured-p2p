package server

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, req *http.Request) {

	fmt.Println("Enviando resposta para o cliente!")

	fmt.Fprintf(w, "Hello from Server!\n")
}

func SearchFlooding(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func SearchRandomWalk(w http.ResponseWriter, req *http.Request) {}

func SearchInDepth(w http.ResponseWriter, req *http.Request) {}

func InitServer(endereco string, porta string) {

	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/SearchFlooding", SearchFlooding)
	http.HandleFunc("/SearchRandomWalk", SearchRandomWalk)
	http.HandleFunc("/SearchInDepth", SearchInDepth)

	fmt.Println("Escutando na porta 10000")
	http.ListenAndServe(":10000", nil)
	http.HandleFunc("/hello", Hello)

	http.ListenAndServe(endereco+":"+porta, nil)
}
