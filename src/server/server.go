package server

import (
	"UP2P/client"
	"UP2P/node"
	"UP2P/utils"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"
)

var NO *node.No

var count int = 4

func Hello(w http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()

	HOST := params.Get("host")
	PORT := params.Get("port")
	NOSEQ := params.Get("noseq")
	TTL := params.Get("ttl")
	MESSAGE_FIELD := params.Get("message")

	MESSAGE := HOST + " " + PORT + " " + NOSEQ + " " + TTL + " " + MESSAGE_FIELD

	fmt.Println("Mensagem recebida: " + MESSAGE)

	node.AddMessage(MESSAGE, NO)

	fmt.Fprintf(w, "")

	tryAddNeighbour(HOST, PORT)

}

func tryAddNeighbour(HOST string, PORT string) {
	vizinho_na_tabela := false

	time.Sleep(100 * time.Millisecond)

	// Checar se o vizinho já está na tabela
	for _, vizinho := range NO.Vizinhos {

		if vizinho.HOST == HOST && vizinho.PORT == PORT {
			vizinho_na_tabela = true
		}

	}

	if !vizinho_na_tabela {
		fmt.Printf("\tAdicionando vizinho na tabela: %s:%s\n", HOST, PORT)

		novo_vizinho := new(node.Vizinho)
		novo_vizinho.HOST = HOST
		novo_vizinho.PORT = PORT

		NO.Vizinhos = append(NO.Vizinhos, novo_vizinho)

		return

	}

	fmt.Printf("\tVizinho ja esta na tabela: %s:%s\n", HOST, PORT)
}

func Search(w http.ResponseWriter, req *http.Request) {

	message_received := utils.ExtrairParamsURL(req)

	fmt.Printf("Mensagem recebida: \"%s\"\n", utils.GenerateStringSearchMessage(message_received))

	is_message_repeated := node.FindReceivedMessage(utils.GenerateStringSearchMessage(message_received), NO)

	if !is_message_repeated {
		node.AddMessage(utils.GenerateStringSearchMessage(message_received), NO)
	} else if message_received.MODE == "FL" {
		fmt.Printf("\tFlooding: mensagem repetida!\n")
		return
	}

	message_to_send := utils.AtualizarMensagemDeBusca(message_received, NO.HOST, NO.PORT)

	value, existsLocally := NO.Pares_chave_valor[message_received.KEY]

	if existsLocally {
		fmt.Printf("\tChave encontrada!\n")
		url := utils.GerarURLdeDevolucao(message_received, value, NO)
		defer http.Get(url)
		return
	}

	//Se TTL iguala a zero a mensagem para aqui
	if message_to_send.TTL == "0" {
		fmt.Println("\tTTL igual a zero, descartando mensagem")
		return
	}

	switch message_received.MODE {
	case "FL":
		SearchFlooding(message_received, message_to_send)
	case "RW":
		SearchRandomWalk(message_received, message_to_send)
	case "BP":
		SearchInDepth(message_received, message_to_send)
	}

}

func SearchFlooding(message_received *utils.SearchMessage, message_to_send *utils.SearchMessage) {

	vizinhos := node.RemoveNeighbour(
		message_received.LAST_HOP_HOST,
		message_received.LAST_HOP_PORT,
		NO.Vizinhos)

	client.ForwardFlooding(message_to_send, vizinhos, NO)

}

func SearchRandomWalk(message_received *utils.SearchMessage, message_to_send *utils.SearchMessage) {

	var random int

	last_hop := message_received.LAST_HOP_HOST + ":" + message_received.LAST_HOP_PORT

	for {
		random = rand.IntN(len(NO.Vizinhos))
		sorted_neighbour := NO.Vizinhos[random].HOST + ":" + NO.Vizinhos[random].PORT

		if last_hop == sorted_neighbour && len(NO.Vizinhos) > 1 {
			continue
		} else {
			break
		}

	}

	url := utils.GerarURLdeSearch(message_to_send, NO, NO.Vizinhos[random])

	client.Consume_endpoint(url, NO, message_to_send, NO.Vizinhos[random])

}

func SearchInDepth(message_received *utils.SearchMessage, message_to_send *utils.SearchMessage) {

	check_message_received := false

	pos := -1

	//Pra verificar se a mensagem é repetida
	for index, dfs_message := range NO.Dfs_messages {

		arr := strings.Split(dfs_message.Message, " ")

		HOST := strings.Split(arr[0], ":")[0]
		PORT := strings.Split(arr[0], ":")[1]

		SEQ_NO := arr[1]

		if message_received.ORIGIN_HOST == HOST &&
			message_received.ORIGIN_PORT == PORT &&
			message_received.NOSEQ == SEQ_NO {

			pos = index
			check_message_received = true
			break

		}

	}

	Last_Hop_Host := message_received.LAST_HOP_HOST
	Last_Hop_Port := message_received.LAST_HOP_PORT

	Received_From := Last_Hop_Host + ":" + Last_Hop_Port

	if check_message_received {

		Active_Child := NO.Dfs_messages[pos].Active_child

		//Pra verificar ciclo

		if Received_From != Active_Child {

			fmt.Printf("BP: ciclo detectado, devolvendo mensagem...\n")

			NO.Dfs_messages[pos].Pending_child =
				node.RemoveNeighbour(Last_Hop_Host, Last_Hop_Port,
					NO.Dfs_messages[pos].Pending_child)

			client.DevolverMensagemDFS(NO.Dfs_messages[pos], NO, Received_From, message_to_send)

			return

		}

		client.SearchInDepth(NO.Dfs_messages[pos], NO, message_to_send)

	} else {

		dfs_message :=
			node.AdicionaMensagemDFS(NO, Received_From, utils.GenerateStringSearchMessage(message_received))

		fmt.Printf("\tDFS Message com %s adicionada\n", dfs_message.Message)

		dfs_message.Pending_child =
			node.RemoveNeighbour(Last_Hop_Host, Last_Hop_Port,
				dfs_message.Pending_child)

		client.SearchInDepth(dfs_message, NO, message_to_send)

	}

}

func KeyReceptor(w http.ResponseWriter, req *http.Request) {
	params := utils.ExtrairParamsURL(req)

	message := params.ORIGIN_HOST + ":" + params.ORIGIN_PORT + " " + params.NOSEQ + " " +
		params.TTL + " " + "VAL" + " " + params.MODE + " " +
		params.KEY + " " + params.VALUE + " " + params.HOP_COUNT

	node.AddMessage(message, NO)

	fmt.Printf("\tValor encontrado!\n\t\tChave: %s valor: %s\n", params.KEY, params.VALUE)

}

func Bye(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	HOST := params.Get("host")
	PORT := params.Get("port")
	NOSEQ := params.Get("noseq")
	TTL := params.Get("ttl")
	MESSAGE_FIELD := params.Get("message")

	MESSAGE := HOST + " " + PORT + " " + NOSEQ + " " + TTL + " " + MESSAGE_FIELD

	fmt.Println("Mensagem recebida: \"" + MESSAGE + "\"")

	node.AddMessage(MESSAGE, NO)

	node.PrintNode(NO, count)

	NO.Vizinhos = node.RemoveNeighbour(HOST, PORT, NO.Vizinhos)
}

func InitServer(_NO *node.No) {

	http.HandleFunc("/Hello", Hello)
	http.HandleFunc("/Search", Search)
	http.HandleFunc("/KeyReceptor", KeyReceptor)
	http.HandleFunc("/Bye", Bye)

	NO = _NO

	fmt.Printf("Servidor criado: %s:%s\n", _NO.HOST, _NO.PORT)
	http.ListenAndServe(_NO.HOST+":"+_NO.PORT, nil)

}
