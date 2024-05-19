package client

import (
	"bufio"
	"fmt"
	"net/http"
)

// Definir funções que o cliente pode executar

func find_all_neighbours() [][]string {

	/* TODO: implementar a busca por todos os vizinhos
	   Returned values: [][]string => matriz de strings com os vizinhos
	*/

	return make([][]string, 0)

}

func search_flooding(_key_ string) (_status_ bool, _value_ int) {
	/* TODO: implementar a requisição da chave para todos os vizinhos
	   Returned values: _status_ => se achou o não a chave
						_key_ => valor da chave, -1 se não for encontrada
	*/

	return true, 1

}

func InitClient(endereco string, porta string) {

	resp, err := http.Get("http://" + endereco + ":" + porta + "/hello")
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
