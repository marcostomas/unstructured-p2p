package client

import (
	"UP2P/node"
	"UP2P/utils"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
)

type SearchMethod func(string, *node.No, string, []*node.Vizinho)

func ShowNeighbours(no *node.No) {

	vizinhos := no.Vizinhos

	fmt.Printf("\nEscolha o vizinho:\n")
	fmt.Printf("Há %d vizinhos na tabela\n", len(vizinhos))

	for i, vizinho := range vizinhos {
		fmt.Printf("\t[%d] %s %s\n", i, vizinho.HOST, vizinho.PORT)
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

	fmt.Println("Encaminhando mensagem \"" +
		no.HOST + ":" +
		no.PORT + " " +
		strconv.Itoa(no.TTL) + " " +
		"1" + " " +
		"HELLO\"" + " para " +
		DESTINY_HOST + ":" + DESTINY_PORT)

	_, err := http.Get(url)

	if err != nil {
		fmt.Printf("\tErro ao conectar!\n")
		return false
	}

	node.IncrementNoSeq(no)

	fmt.Println("\tEnvio feito com sucesso: " +
		no.HOST + ":" +
		no.PORT + " " +
		noseq + " " +
		"1")

	return true
}

func FindKey(no *node.No, f SearchMethod, KEY_AUTO string) {

	fmt.Printf("Digite a chave a ser buscada\n")

	var KEY string

	if KEY_AUTO == "NONE" {
		fmt.Scanln(&KEY)
	} else {
		KEY = KEY_AUTO
	}

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

	for index := range Vizinhos {
		url := utils.GerarURLdeSearch(message, NO, NO.Vizinhos[index])
		Consume_endpoint(url, NO, message, NO.Vizinhos[index])
	}
}

func ForwardFlooding(MESSAGE *utils.SearchMessage, Vizinhos []*node.Vizinho, NO *node.No) {

	for index := range Vizinhos {
		url := utils.GerarURLdeSearch(MESSAGE, NO, Vizinhos[index])
		go Consume_endpoint(url, NO, MESSAGE, Vizinhos[index])
	}

}

func SearchRandomWalk(KEY string, NO *node.No, TTL string, Vizinhos []*node.Vizinho) {
	random := rand.IntN(len(NO.Vizinhos))

	message := utils.GerarMensagemDeBusca(NO, TTL, "RW", KEY)

	url := utils.GerarURLdeSearch(message, NO, NO.Vizinhos[random])

	Consume_endpoint(url, NO, message, NO.Vizinhos[random])

	node.IncrementNoSeq(NO)
}

func PrepareSearchInDepth(KEY string,
	NO *node.No,
	TTL string,
	VIZINHOS []*node.Vizinho) {

	message := utils.GenerateStringSearchMessage(
		utils.GerarMensagemDeBusca(NO, TTL, "BP", KEY))

	node_address := NO.HOST + ":" + NO.PORT

	dfs_message := node.AdicionaMensagemDFS(NO, node_address, message)

	message_to_send := utils.ConverterDFSMessage(dfs_message, "")

	SearchInDepth(dfs_message, NO, message_to_send)

	node.IncrementNoSeq(NO)

}

func SearchInDepth(DFS_MESSAGE *node.DfsMessage,
	NO *node.No,
	message_to_send *utils.SearchMessage) {

	//Checa se tem vizinhos pendentes na dfs_message do nó
	if len(DFS_MESSAGE.Pending_child) == 0 {

		fmt.Printf("BP: nenhum vizinho encontrou a chave, retrocedendo...\n")

		//Checa se o nó que enviou a mensagem não é o mesmo nó
		if DFS_MESSAGE.Received_from != NO.HOST+":"+
			NO.PORT {

			pos := 0

			//Precisamos da posição do vizinho que enviou a mensagem
			for index, vizinho := range NO.Vizinhos {
				pos = index
				if vizinho.HOST+":"+vizinho.PORT == DFS_MESSAGE.Received_from {
					break
				}
			}

			message := utils.ConverterDFSMessage(DFS_MESSAGE, "")

			vizinho := NO.Vizinhos[pos]

			url := utils.GerarURLdeSearch(message, NO, vizinho)

			Consume_endpoint(url, NO, message, vizinho)

			return

		} else {
			KEY := strings.Split(DFS_MESSAGE.Message, " ")[6]
			fmt.Printf("BP: Não foi possível localizar a chave %s\n", KEY)
			return
		}

	}

	neighbour := utils.EscolherVizinhoAleatorio(DFS_MESSAGE)

	DFS_MESSAGE.Active_child = neighbour.HOST + ":" + neighbour.PORT

	url := utils.GerarURLdeSearch(message_to_send, NO, neighbour)

	Consume_endpoint(url, NO, message_to_send, neighbour)

}

func DevolverMensagemDFS(DFS_MESSAGE *node.DfsMessage, NO *node.No, RECEIVED_FROM string,
	MESSAGE_TO_SEND *utils.SearchMessage) {

	pos := 0

	//Precisamos da posição do vizinho que enviou a mensagem
	for index, vizinho := range NO.Vizinhos {
		pos = index
		if vizinho.HOST+":"+vizinho.PORT == RECEIVED_FROM {
			break
		}
	}

	vizinho := NO.Vizinhos[pos]

	url := utils.GerarURLdeSearch(MESSAGE_TO_SEND, NO, vizinho)

	Consume_endpoint(url, NO, MESSAGE_TO_SEND, vizinho)
}

func Bye(NO *node.No) {
	for _, vizinho := range NO.Vizinhos {

		message := NO.HOST + ":" + NO.PORT +
			" " + strconv.Itoa(NO.NoSeq) + " " + "1" + " " + "BYE"

		fmt.Printf("Encaminhando mensagem \"%s\" para %s:%s\n", message, vizinho.HOST, vizinho.PORT)

		var url string = "http://" +
			vizinho.HOST +
			":" +
			vizinho.PORT +
			"/Bye?" +
			"host=" + NO.HOST +
			"&port=" + NO.PORT +
			"&noseq=" + strconv.Itoa(NO.NoSeq) +
			"&ttl=" + "1" +
			"&message=" + "BYE"

		http.Get(url)

		fmt.Printf("\tEnvio feito com sucesso: \"%s\"", message)

	}
}

func Consume_endpoint(url string, no *node.No, MESSAGE *utils.SearchMessage,
	Vizinho *node.Vizinho) bool {

	fmt.Println("Encaminhando mensagem \"" +
		MESSAGE.ORIGIN_HOST + ":" +
		MESSAGE.ORIGIN_PORT + " " +
		MESSAGE.NOSEQ + " " +
		MESSAGE.TTL + " " +
		MESSAGE.ACTION + " " +
		MESSAGE.MODE + " " +
		MESSAGE.HOP_COUNT + "\"" +
		" para " + Vizinho.HOST + ":" + Vizinho.PORT)

	_, err := http.Get(url)

	if err != nil {
		fmt.Printf("\tErro ao conectar!\n")
		return false
	}

	fmt.Println("\tEnvio feito com sucesso: \"" +
		MESSAGE.ORIGIN_HOST + ":" +
		MESSAGE.ORIGIN_PORT + " " +
		MESSAGE.NOSEQ + " " +
		MESSAGE.TTL + " " +
		MESSAGE.ACTION + " " +
		MESSAGE.MODE + " " +
		MESSAGE.LAST_HOP_PORT + " " +
		MESSAGE.KEY + " " +
		MESSAGE.HOP_COUNT + "\"")

	return true
}
