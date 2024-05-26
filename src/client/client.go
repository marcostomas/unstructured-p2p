package client

import (
	"UP2P/node"
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Definir funções que o cliente pode executar

func ListAllNeighbours() [][]string {

	/* TODO: implementar a busca por todos os vizinhos
	   Returned values: [][]string => matriz de strings com os vizinhos
	*/

	return make([][]string, 0)

}

func ShowNeighboursToChoose(no *node.No) {

	vizinhos := no.Vizinhos

	fmt.Printf("\nEscolha o vizinho:\n")
	fmt.Printf("Há %d vizinhos na tabela\n", len(vizinhos))

	for i, vizinho := range vizinhos {
		fmt.Printf("[%d] %s:%s\n", i, vizinho.HOST, vizinho.PORT)
	}

	var n int
	_, err := fmt.Scanln(&n)

	if err != nil {
		fmt.Println("Erro ao ler o número", err)
	}

	Hello(
		vizinhos[n].HOST,
		vizinhos[n].PORT,
		no)

}

func Hello(DESTINY_HOST string,
	DESTINY_PORT string,
	no *node.No) bool {

	//Converter de int para string
	noseq := strconv.Itoa(no.NoSeq)

	var url string = "http://" +
		DESTINY_HOST +
		":" +
		DESTINY_PORT +
		"/Hello?" +
		"host=" + no.HOST +
		"&port=" + no.PORT +
		"&noseq=" + noseq +
		"&ttl=" + "1" +
		"&message=" + "HELLO"

	node.IncrementNoSeq(no)

	fmt.Println("Encaminhando mensagem \"" +
		no.HOST + ":" +
		no.PORT + " " +
		noseq + " " +
		"1" + " " +
		"HELLO" + "\"" +
		" para " + DESTINY_HOST + ":" + DESTINY_PORT)

	message, status := consumeEndpoint(url)

	if !status {
		fmt.Println("Não foi possível fazer a comunicação com: " + DESTINY_HOST + ":" + DESTINY_PORT)
		fmt.Println("Motivo: " + message)
		return false
	}

	fmt.Println("\tEnvio feito com sucesso: " +
		no.HOST + ":" +
		no.PORT + " " +
		noseq + " " +
		"1")
	return true

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

	if err != nil {
		return "Não foi possível estabelecer a conexão com " + url, false
	}

	if resp.StatusCode == 404 {
		return "404, recurso não encontrado", false
	}

	var message string

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		message += scanner.Text()
	}

	return message, true
}
