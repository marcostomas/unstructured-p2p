package utils

import (
	"UP2P/node"
	"net/http"
	"strconv"
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
	posVizinho int) string {

	// Converter de int para string
	noseq := strconv.Itoa(NO.NoSeq)

	return "http://" +
		NO.Vizinhos[posVizinho].HOST + ":" +
		NO.Vizinhos[posVizinho].PORT + "/" +
		"Search" + "?" +
		"host=" + message.ORIGIN_HOST + "&" +
		"port=" + message.ORIGIN_PORT + "&" +
		"seqno=" + noseq + "&" +
		"ttl=" + message.TTL + "&" +
		"action=" + "search" + "&" +
		"mode=" + message.MODE + "&" +
		"last_hop_port=" + NO.PORT + "&" +
		"key=" + message.KEY + "&" +
		"hop_count=" + "1"
}

func GerarURLdeDevolucao(
	message *SearchMessage,
	VALUE string,
	NO *node.No) string {

	// Converter de int para string
	noseq := strconv.Itoa(NO.NoSeq)

	return "http://" +
		message.ORIGIN_HOST + ":" +
		message.ORIGIN_PORT + "/" +
		"KeyReceptor" + "?" +
		"host=" + NO.HOST + "&" +
		"port=" + NO.PORT + "&" +
		"seqno=" + noseq + "&" +
		"ttl=" + message.TTL + "&" +
		"action=" + "VAL" + "&" +
		"mode=" + message.MODE + "&" +
		"key=" + message.KEY + "&" +
		"value=" + VALUE + "&" +
		"hop_count=" + message.HOP_COUNT
}
