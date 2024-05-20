package server

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, req *http.Request) {

	fmt.Println("\n#############################################")
	fmt.Println("###### Enviando resposta para o cliente! ####")
	fmt.Println("#############################################\n")

	fmt.Fprintf(w, "Hello from Server!\n")

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

func InitServer(PORT string) {

	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/SearchFlooding", SearchFlooding)
	http.HandleFunc("/SearchRandomWalk", SearchRandomWalk)
	http.HandleFunc("/SearchInDepth", SearchInDepth)

	fmt.Printf("Escutando na porta %s\n", PORT)
	http.ListenAndServe(PORT, nil)
}
