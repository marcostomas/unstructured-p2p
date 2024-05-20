package client

import (
	"bufio"
	"fmt"
	"net/http"
)

// Definir funções que o cliente pode executar

func list_all_neighbours(_file_ string) {

	/* TODO: imprimir todos os vizinhos
	   Returned values: void
	*/

}

func search_flooding(_key_ string) (_status_ bool, _value_ int) {
	/* TODO: implementar a requisição da chave para todos os vizinhos usando flooding
	   Returned values: _status_ => se achou o não a chave
						_key_ => valor da chave, -1 se não for encontrada
	*/

	return true, 1

}

func search_random_walk(_key_ string) (_status_ bool, _value_ int) {
	/*  TODO: implementar a requisição de chave para todos os veizinhos usando random walk
	   Returned values: _status_ => se achou o não a chave
					_key_ => valor da chave, -1 se não for encontrada
	*/

	return true, 1
}

func init_client() {

	resp, err := http.Get("http://localhost:8090/hello")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
