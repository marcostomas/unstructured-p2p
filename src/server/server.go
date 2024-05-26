package server

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()

	HOST := params.Get("host")
	PORT := params.Get("port")
	NOSEQ := params.Get("noseq")
	TTL := params.Get("ttl")
	MESSAGE := params.Get("message")

	fmt.Println("Mensagem recebida: " +
		HOST + ":" +
		PORT + " " +
		NOSEQ + " " +
		TTL + " " +
		MESSAGE)

	fmt.Fprintf(w, "OK!")

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

func InitServer(HOST string, PORT string) {

	http.HandleFunc("/Hello", Hello)
	http.HandleFunc("/SearchFlooding", SearchFlooding)
	http.HandleFunc("/SearchRandomWalk", SearchRandomWalk)
	http.HandleFunc("/SearchInDepth", SearchInDepth)

	fmt.Printf("Escutando na porta %s\n", PORT)
	http.ListenAndServe(HOST+":"+PORT, nil)
}
