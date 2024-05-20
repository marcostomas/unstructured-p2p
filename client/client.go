package client

import (
	"bufio"
	"fmt"
	"net/http"
	"time"
)

// Definir funções que o cliente pode executar

func ListAllNeighbours() [][]string {

	/* TODO: implementar a busca por todos os vizinhos
	   Returned values: [][]string => matriz de strings com os vizinhos
	*/

	return make([][]string, 0)

}

func ShowNeighboursToChoose(vizinhos []string, handler func(string) bool) {

	fmt.Printf("\nEscolha o vizinho:\n")
	fmt.Printf("Há %d vizinhos na tabela\n", len(vizinhos))

	for i, vizinho := range vizinhos {
		fmt.Printf("[%d] %s\n", i, vizinho)
	}

	var numero int
	_, err := fmt.Scanln(&numero)

	if err != nil {
		fmt.Println("Erro ao ler o número", err)
	}

	handler(vizinhos[numero])

}

func Hello(endereco string) bool {

	var url string = "http://" + endereco + "/hello"

	message, status := consumeEndpoint(url)

	if status {
		fmt.Println("Message received from " + endereco + ": " + message)
		return true
	} else {
		fmt.Println("Não foi possível fazer a comunicação com " + endereco)
		return false
	}

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

func consumeEndpoint(url string) (string, bool) {

	time.Sleep(1000 * time.Millisecond)

	resp, err := http.Get(url)

	var message string

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		message += scanner.Text()
	}

	if err != nil {
		return message, false
	}

	return message, true
}
