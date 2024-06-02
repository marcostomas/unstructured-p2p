package server

import (
	"UP2P/node"
	"UP2P/utils"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
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
		fmt.Printf(err.Error())
	}

	fmt.Printf("\nMensagem recebida: \n%s\n", string(data))

	//Cast de string pra int das propriedades TTL e HOP_COUNT
	TTL, _ := strconv.Atoi(message.TTL)
	HOP_COUNT, _ := strconv.Atoi(message.HOP_COUNT)

	TTL--
	HOP_COUNT++

	//Se TTL iguala a zero a mensagem para aqui
	if TTL == 0 {
		return
	}

	//Modifica a mensagem para o próximo envio
	message.TTL = strconv.Itoa(TTL)
	message.HOP_COUNT = strconv.Itoa(HOP_COUNT)

	switch message.MODE {
	case "FL":
		SearchFlooding(message)
		break
	case "RW":
		SearchRandomWalk(message)
		break
	case "DP":
		SearchInDepth(message)
		break
	}

}

func SearchFlooding(message *utils.SearchMessage) {

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

	fmt.Printf("Escutando na porta %s\n", _NO.PORT)
	http.ListenAndServe(_NO.HOST+":"+_NO.PORT, nil)

}
