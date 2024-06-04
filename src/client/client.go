package client

import (
	"UP2P/node"
	"UP2P/utils"
	"bufio"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

type SearchMethod func(string, *node.No, string, []*node.Vizinho)

// Definir funções que o cliente pode executar

func ListAllNeighbours() [][]string {

	/* TODO: implementar a busca por todos os vizinhos
	   Returned values: [][]string => matriz de strings com os vizinhos
	*/

	return make([][]string, 0)

}

func ShowNeighbours(no *node.No) {

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

func FindKey(no *node.No, f SearchMethod) {

	fmt.Printf("Digite a chave a ser buscada\n")

	var KEY string

	fmt.Scanln(&KEY)

	value, existsLocally := no.Pares_chave_valor[KEY]

	if existsLocally {
		fmt.Printf("Valor na tabela local!\n")
		fmt.Printf("\tchave: %s valor: %s\n", KEY, value)
		return
	}

	f(KEY, no, strconv.Itoa(no.TTL), no.Vizinhos)

}

func SearchFlooding(KEY string,
	NO *node.No,
	TTL string,
	Vizinhos []*node.Vizinho) {

	message := utils.GerarMensagemDeBusca(NO, TTL, "FL", KEY)
	sMsg := fmt.Sprintf("%s:%s %s %s %s %s %s %s",
		message.ORIGIN_HOST,
		message.ORIGIN_PORT,
		message.NOSEQ,
		message.TTL,
		message.ACTION,
		message.MODE,
		message.LAST_HOP_PORT,
		message.KEY)

	node.AddMessage(sMsg, NO)
	node.IncrementNoSeq(NO)

	for index, _ := range Vizinhos {
		url := utils.GerarURLdeSearch(message, NO, index)
		http.Get(url)
	}
}

func SearchRandomWalk(KEY string, NO *node.No, TTL string, Vizinhos []*node.Vizinho) {
	random := rand.IntN(len(NO.Vizinhos))

	message := utils.GerarMensagemDeBusca(NO, TTL, "RW", KEY)

	url := utils.GerarURLdeSearch(message, NO, random)

	http.Get(url)

	node.IncrementNoSeq(NO)
}

func SearchInDepth(KEY string, NO *node.No, TTL string, Vizinhos []*node.Vizinho) {

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
