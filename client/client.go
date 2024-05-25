package client

import (
	"fmt"
	"net/http"
)

// Definir funções que o cliente pode executar

func ListAllNeighbours() [][]string {
	/* TODO: implementar a busca por todos os vizinhos
	   Returned values: [][]string => matriz de strings com os vizinhos
	*/
	return make([][]string, 0)
}

func Hello() {
	consumeEndpoint("http://localhost:10000/hello")
}

func SearchFlooding(_key_ string) (_status_ bool, _value_ int) {
	/* TODO: implementar a requisição da chave para todos os vizinhos
	   Returned values: _status_ => se achou o não a chave
						_key_ => valor da chave, -1 se não for encontrada
	*/

	fmt.Println("Mandando um searchFlooding")
	return true, 1
}

func SearchRandomWalk(_key_ string) (_status_ bool, _value_ int) {
	return true, 1
}

func SearchInDepth(_key_ string) (_status_ bool, _value_ int) {
	return true, 1
}

func bye(uri string) bool {
	resp, err := http.Get("http://" + uri + "/bye")
	if err != nil {
		fmt.Println("Erro ao dizer xau:", err)
		return false
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
	return true
}

func consumeEndpoint(url string) {
	// resp, err := http.Get(url)
}
