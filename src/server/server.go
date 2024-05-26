package server

import (
	"fmt"
	"net/http"
)

type vizinho struct {
	HOST string
	PORT string
}

type no struct {
	HOST              string
	PORT              string
	seqNo             int
	pares_chave_valor map[string]string //Par nome e número associado
	vizinhos          []*vizinho
}

func Hello(w http.ResponseWriter, req *http.Request) {

	url := req.URL.Query()

	HOST := url.Get("host")
	PORT := url.Get("port")
	NOSEQ := url.Get("ttl")
	TTL := url.Get("ttl")
	MESSAGE := url.Get("message")

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

	http.HandleFunc("/hello/?host=host&port=port&noseq=noseq&ttl=ttl&message=message", Hello)
	http.HandleFunc("/SearchFlooding", SearchFlooding)
	http.HandleFunc("/SearchRandomWalk", SearchRandomWalk)
	http.HandleFunc("/SearchInDepth", SearchInDepth)

	fmt.Printf("Escutando na porta %s\n", PORT)
	http.ListenAndServe(HOST+":"+PORT, nil)
}
