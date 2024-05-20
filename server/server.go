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

<<<<<<< HEAD
func init_server() {
=======
func SearchRandomWalk(w http.ResponseWriter, req *http.Request) {
>>>>>>> 1c6ac3a3466a4e05bd8c4c98b574149fc2d0fccf

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
