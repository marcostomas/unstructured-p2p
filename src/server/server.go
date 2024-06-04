package server

import (
	"UP2P/client"
	"UP2P/node"
	"UP2P/utils"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
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

	message := utils.ExtrairParamsURL(req)

	data, err := json.Marshal(message)

	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Printf("\nMensagem recebida: \n%s\n", string(data))

	new_message := utils.AtualizarMensagemDeBusca(message, NO.PORT)

	TTL, _ := strconv.Atoi(message.TTL)

	//Se TTL iguala a zero a mensagem para aqui
	if TTL == 0 {
		fmt.Println("TTL igual a zero, descartando mensagem")
		return
	}

	req_host := strings.Split(req.RemoteAddr, ":")[0]

	switch message.MODE {
	case "FL":
		SearchFlooding(message, new_message, req_host)
	case "RW":
		SearchRandomWalk(message)
	case "DP":
		SearchInDepth(message)
	}

}

func SearchFlooding(message *utils.SearchMessage, new_message *utils.SearchMessage, req_host string) {
	msg_received := node.FindReceivedMessage(utils.GenerateStringSearchMessage(message), NO)

	if msg_received {
		fmt.Println("Mensagem já recebida!")
		return
	}

	node.AddMessage(utils.GenerateStringSearchMessage(message), NO)

	value, exists := NO.Pares_chave_valor[message.KEY]

	if exists {
		fmt.Println("Chave encontrada!")
		return_url := utils.GerarURLdeDevolucao(message, value, NO)
		defer http.Get(return_url)
		node.IncrementNoSeq(NO)
		return
	}

	fmt.Printf("A chave %s não foi encontrada na tabela local!", message.KEY)

	client.ForwardFlooding(new_message, node.RemoveNeighbour(req_host, message.LAST_HOP_PORT, NO),
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

	url := utils.GerarURLdeSearch(message, NO, random)

	fmt.Printf("URL gerada para novo envio: %s\n", url)

	fmt.Println("Encaminhando mensagem para " + NO.Vizinhos[random].HOST + ":" + NO.Vizinhos[random].PORT)
	defer http.Get(url)
	node.IncrementNoSeq(NO)
	fmt.Printf("NoSeq incrementando: %d\n", NO.NoSeq)

}

func SearchInDepth(message *utils.SearchMessage) {

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
