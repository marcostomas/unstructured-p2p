package server

import (
	"UP2P/client"
	"UP2P/node"
	"UP2P/utils"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
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

	node.PrintNode(NO, count)

	count++

	fmt.Fprintf(w, "OK!")

}

func Search(w http.ResponseWriter, req *http.Request) {

	message_received := utils.ExtrairParamsURL(req)

	message_updated := utils.AtualizarMensagemDeBusca(message_received, NO.PORT)

	//Se TTL iguala a zero a mensagem para aqui
	if message_updated.TTL == "0" {
		fmt.Println("TTL igual a zero, descartando mensagem")
		return
	}

	req_host := strings.Split(req.RemoteAddr, ":")[0]

	switch message_received.MODE {
	case "FL":
		SearchFlooding(message_received, message_updated, req_host)
	case "RW":
		SearchRandomWalk(message_received)
	case "DP":
		SearchInDepth(message_received, message_updated)
	}

}

func SearchFlooding(message_received *utils.SearchMessage, message_updated *utils.SearchMessage, req_host string) {

	msg_received := node.FindReceivedMessage(utils.GenerateStringSearchMessage(message_received), NO)

	if msg_received {
		fmt.Println("Mensagem já recebida!")
		return
	}

	node.AddMessage(utils.GenerateStringSearchMessage(message_received), NO)

	value, exists := NO.Pares_chave_valor[message_received.KEY]

	if exists {
		fmt.Println("Chave encontrada!")
		return_url := utils.GerarURLdeDevolucao(message_received, value, NO)
		defer http.Get(return_url)
		node.IncrementNoSeq(NO)
		return
	}

	fmt.Printf("A chave %s não foi encontrada na tabela local!", message_received.KEY)

	client.ForwardFlooding(message_updated, node.RemoveNeighbour(req_host, message_updated.LAST_HOP_PORT, NO.Vizinhos),
		NO)

}

func SearchRandomWalk(message *utils.SearchMessage) {
	value, exists := NO.Pares_chave_valor[message.KEY]

	if exists {
		fmt.Printf("\nValor da chave %s encontrado: %s!\n", message.KEY, value)
		url := utils.GerarURLdeDevolucao(message, value, NO)
		fmt.Printf("%s\n", url)
		defer http.Get(url)
		node.IncrementNoSeq(NO)
		fmt.Printf("NoSeq incrementando: %d\n", NO.NoSeq)
		return
	}

	fmt.Printf("A chave %s não foi encontrada na tabela local!", message.KEY)

	random := rand.IntN(len(NO.Vizinhos))

	url := utils.GerarURLdeSearch(message, NO, NO.Vizinhos[random])

	fmt.Printf("URL gerada para novo envio: %s\n", url)

	fmt.Println("Encaminhando mensagem para " + NO.Vizinhos[random].HOST + ":" + NO.Vizinhos[random].PORT)
	defer http.Get(url)
	node.IncrementNoSeq(NO)
	fmt.Printf("NoSeq incrementando: %d\n", NO.NoSeq)

}

func SearchInDepth(message_received *utils.SearchMessage, message_updated *utils.SearchMessage) {

	fmt.Printf("Search in Depth-----------------\n")
	fmt.Printf("Mensagem recebida %s\n", utils.GenerateStringSearchMessage(message_received))

	check_message_received := false

	pos := 0

	//Pra verificar se a mensagem é repitida
	for index, dfs_message := range NO.Dfs_messages {

		fmt.Printf("PÉREZ!\n")
		fmt.Printf("Tá barato!\n")

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

	Host_Active_Child := strings.Split(NO.Dfs_messages[pos].Active_child, ":")[0]
	Port_Active_Child := strings.Split(NO.Dfs_messages[pos].Active_child, ":")[1]

	if check_message_received {

		//Pra verificar ciclo
		if message_received.LAST_HOP_PORT != Port_Active_Child {

			fmt.Printf("BP: ciclo detectado, devolvendo mensagem...")

			NO.Dfs_messages[pos].Pending_child =
				node.RemoveNeighbour(Host_Active_Child, Port_Active_Child,
					NO.Dfs_messages[pos].Pending_child)

		}

		client.SearchInDepth(NO.Dfs_messages[pos], NO)

	} else {

		fmt.Printf("Mensagem nova!\n")

		dfs_message :=
			utils.AdicionaMensagemDFS(NO, utils.GenerateStringSearchMessage(message_updated))

		fmt.Printf("DFS Message com %s adicionada\n", dfs_message.Message)

		dfs_message.Pending_child =
			node.RemoveNeighbour(Host_Active_Child, Port_Active_Child,
				dfs_message.Pending_child)

		client.SearchInDepth(dfs_message, NO)

		fmt.Printf("PÉREZ\n")
		fmt.Printf("--------------------\n")

	}

}

func KeyReceptor(w http.ResponseWriter, req *http.Request) {
	params := utils.ExtrairParamsURL(req)

	message := params.ORIGIN_HOST + ":" + params.ORIGIN_PORT + " " + params.NOSEQ + " " +
		params.TTL + " " + "VAL" + " " + params.MODE + " " +
		params.KEY + " " + params.VALUE + " " + params.HOP_COUNT

	node.AddMessage(message, NO)

	fmt.Printf("Valor encontrado! Chave: %s valor: %s\n", params.KEY, params.VALUE)

}

func InitServer(_NO *node.No) {

	http.HandleFunc("/Hello", Hello)
	http.HandleFunc("/Search", Search)
	http.HandleFunc("/KeyReceptor", KeyReceptor)

	fmt.Println(_NO)

	NO = _NO

	fmt.Println(NO)

	fmt.Printf("Escutando na porta %s:%s\n", _NO.HOST, _NO.PORT)
	http.ListenAndServe(_NO.HOST+":"+_NO.PORT, nil)

}
