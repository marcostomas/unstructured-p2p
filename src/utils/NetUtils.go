package utils

import (
	"UP2P/node"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
)

type SearchMessage struct {
	ORIGIN_HOST   string
	ORIGIN_PORT   string
	NOSEQ         string
	TTL           string
	ACTION        string
	MODE          string
	LAST_HOP_PORT string
	KEY           string
	VALUE         string
	HOP_COUNT     string
}

func GerarMensagemDeBusca(NO *node.No, _TTL string, _MODE string, _KEY string) *SearchMessage {

	// Converter de int para string
	noseq := strconv.Itoa(NO.NoSeq)

	return &SearchMessage{
		ORIGIN_HOST:   NO.HOST,
		ORIGIN_PORT:   NO.PORT,
		NOSEQ:         noseq,
		TTL:           _TTL,
		ACTION:        "SEARCH",
		MODE:          _MODE,
		LAST_HOP_PORT: NO.PORT,
		KEY:           _KEY,
		VALUE:         "",
		HOP_COUNT:     "1",
	}
}

func AtualizarMensagemDeBusca(MESSAGE *SearchMessage, PORT string) *SearchMessage {

	TTL, _ := strconv.Atoi(MESSAGE.TTL)
	TTL--

	HOP_COUNT, _ := strconv.Atoi(MESSAGE.HOP_COUNT)
	HOP_COUNT++

	MESSAGE.TTL = strconv.Itoa(TTL)
	MESSAGE.HOP_COUNT = strconv.Itoa(HOP_COUNT)

	MESSAGE.LAST_HOP_PORT = PORT

	return &SearchMessage{
		ORIGIN_HOST:   MESSAGE.ORIGIN_HOST,
		ORIGIN_PORT:   MESSAGE.ORIGIN_PORT,
		NOSEQ:         MESSAGE.NOSEQ,
		TTL:           strconv.Itoa(TTL),
		ACTION:        MESSAGE.ACTION,
		MODE:          MESSAGE.MODE,
		LAST_HOP_PORT: PORT,
		KEY:           MESSAGE.KEY,
		VALUE:         "",
		HOP_COUNT:     strconv.Itoa(HOP_COUNT),
	}

}

func ConverterDFSMessage(DFS_MESSAGE *node.DfsMessage, VALUE string) *SearchMessage {
	arr := strings.Split(DFS_MESSAGE.Message, " ")

	ORIGIN_HOST := strings.Split(arr[0], ":")[0]
	ORIGIN_PORT := strings.Split(arr[0], ":")[1]

	return &SearchMessage{
		ORIGIN_HOST:   ORIGIN_HOST,
		ORIGIN_PORT:   ORIGIN_PORT,
		NOSEQ:         arr[1],
		TTL:           arr[2],
		ACTION:        arr[3],
		MODE:          arr[4],
		LAST_HOP_PORT: arr[5],
		KEY:           arr[6],
		VALUE:         VALUE,
		HOP_COUNT:     arr[8],
	}

}

// Retorna um ponteiro para uma struct SearchMessage, definida no começo desse arquivo,
// com os parâmetros da url
func ExtrairParamsURL(req *http.Request) *SearchMessage {
	params := req.URL.Query()

	_HOST := params.Get("host")
	_PORT := params.Get("port")
	_NOSEQ := params.Get("seqno")
	_TTL := params.Get("ttl")
	_ACTION := params.Get("action")
	_MODE := params.Get("mode")
	_LAST_HOP_PORT := params.Get("last_hop_port")
	_VALUE := params.Get("value")
	_KEY := params.Get("key")
	_HOP_COUNT := params.Get("hop_count")

	message := &SearchMessage{
		ORIGIN_HOST:   _HOST,
		ORIGIN_PORT:   _PORT,
		NOSEQ:         _NOSEQ,
		TTL:           _TTL,
		ACTION:        _ACTION,
		MODE:          _MODE,
		LAST_HOP_PORT: _LAST_HOP_PORT,
		KEY:           _KEY,
		VALUE:         _VALUE,
		HOP_COUNT:     _HOP_COUNT,
	}

	return message

}

func GerarURLdeSearch(
	message *SearchMessage,
	NO *node.No,
	VIZINHO *node.Vizinho) string {

	return "http://" +
		VIZINHO.HOST + ":" +
		VIZINHO.PORT + "/" +
		"Search" + "?" +
		"host=" + message.ORIGIN_HOST + "&" +
		"port=" + message.ORIGIN_PORT + "&" +
		"seqno=" + message.NOSEQ + "&" +
		"ttl=" + message.TTL + "&" +
		"action=" + "search" + "&" +
		"mode=" + message.MODE + "&" +
		"last_hop_port=" + NO.PORT + "&" +
		"key=" + message.KEY + "&" +
		"hop_count=" + message.HOP_COUNT
}

func GerarURLdeDevolucao(
	message *SearchMessage,
	VALUE string,
	NO *node.No) string {

	// Converter de int para string
	noseq := strconv.Itoa(NO.NoSeq)
	ttl := strconv.Itoa(NO.TTL)

	return "http://" +
		message.ORIGIN_HOST + ":" +
		message.ORIGIN_PORT + "/" +
		"KeyReceptor" + "?" +
		"host=" + NO.HOST + "&" +
		"port=" + NO.PORT + "&" +
		"seqno=" + noseq + "&" +
		"ttl=" + ttl + "&" +
		"action=" + "VAL" + "&" +
		"mode=" + message.MODE + "&" +
		"key=" + message.KEY + "&" +
		"value=" + VALUE + "&" +
		"hop_count=" + message.HOP_COUNT
}

func GenerateStringSearchMessage(message *SearchMessage) string {
	return fmt.Sprintf("%s:%s %s %s %s %s %s %s %s %s", message.ORIGIN_HOST,
		message.ORIGIN_PORT, message.NOSEQ, message.TTL, message.ACTION, message.MODE,
		message.LAST_HOP_PORT, message.KEY, message.VALUE, message.HOP_COUNT)
}

// Escolhe um vizinho aleatoriamente e remove esse vizinho dos vizinhos pendentes
func EscolherVizinhoAleatorio(DFS_MESSAGE *node.DfsMessage) *node.Vizinho {
	//Escolhe um vizinho aleatório
	random := rand.IntN(len(DFS_MESSAGE.Pending_child))

	neighbour := DFS_MESSAGE.Pending_child[random]

	DFS_MESSAGE.Pending_child = append(DFS_MESSAGE.Pending_child[:random], DFS_MESSAGE.Pending_child[random+1:]...)

	return neighbour
}

func AdicionaMensagemDFS(noh *node.No, origem_msg string) *node.DfsMessage {
	temp := &node.DfsMessage{
		Message:       origem_msg,
		Received_from: noh.HOST + ":" + noh.PORT,
		Active_child:  "",
		Pending_child: make([]*node.Vizinho, 0),
	}

	temp.Pending_child = append(temp.Pending_child, noh.Vizinhos...)

	noh.Dfs_messages = append(noh.Dfs_messages, temp)
	return temp
}
