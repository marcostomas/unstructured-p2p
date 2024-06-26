package node

import (
	"fmt"
	"strings"
)

type Vizinho struct {
	HOST string
	PORT string
}

type No struct {
	HOST              string
	PORT              string
	NoSeq             int
	Received_messages []string
	Pares_chave_valor map[string]string //Par nome e número associado
	Vizinhos          []*Vizinho
	TTL               int
	Dfs_messages      []*DfsMessage //
}

type DfsMessage struct {
	Message       string     // Mensagem recebida
	Received_from string     // Vizinho que enviou a mensagem
	Active_child  string     // Vizinho que está ativo, ou seja, pra quem ele enviou o pedido de busca
	Pending_child []*Vizinho // Vizinhos que ainda não foram visitados
}

func NewNo(_HOST string, _PORT string) *No {

	return &No{
		HOST:              _HOST,
		PORT:              _PORT,
		NoSeq:             1, // Primeiro envia para depois incrementar
		Pares_chave_valor: map[string]string{},
		Vizinhos:          make([]*Vizinho, 0),
		Received_messages: make([]string, 0),
		TTL:               100,
		Dfs_messages:      make([]*DfsMessage, 0),
	}
}

func AddKey(par string, noh *No) {

	key := strings.Split(par, " ")[0]
	value := strings.Split(par, " ")[1]

	noh.Pares_chave_valor[key] = value
}

func AddNeighbour(no *No, _HOST string, _PORT string) {
	no.Vizinhos = append(no.Vizinhos, &Vizinho{HOST: _HOST, PORT: _PORT})
}

func IncrementNoSeq(no *No) {
	no.NoSeq++
}

func AddMessage(MESSAGE string, NO *No) {
	NO.Received_messages = append(NO.Received_messages, MESSAGE)
}

func GenerateNeighboursList(data []byte) []*Vizinho {
	arr := strings.Split(string(data), "\n")
	listaVizinhos := make([]*Vizinho, 0)
	for _, vizinho := range arr {
		listaVizinhos = append(listaVizinhos, &Vizinho{
			HOST: strings.Split(vizinho, ":")[0],
			PORT: strings.Split(vizinho, ":")[1]})
	}

	return listaVizinhos

}

func PrintNode(noh *No, count int) {

	fmt.Printf("\n\n")

	fmt.Println("PORT: ", noh.PORT)
	fmt.Println("Vizinhos: [")

	for _, vizinho := range noh.Vizinhos {
		fmt.Printf("\t%s:%s\n", vizinho.HOST, vizinho.PORT)
	}

	fmt.Printf("]\n")

	fmt.Printf("\n")
}

func FindReceivedMessage(message string, NO *No) bool {

	paramsMsg := strings.Split(message, " ")

	for _, msg := range NO.Received_messages {

		paramsOtherMsg := strings.Split(msg, " ")

		addressOtherMsg := paramsOtherMsg[0]

		if paramsMsg[0] == addressOtherMsg &&
			paramsMsg[1] == paramsOtherMsg[1] {
			return true
		}
	}
	return false
}

func ChangeTTL(no *No) {
	fmt.Println("Digite o novo valor de TTL")
	fmt.Scanln(&no.TTL)
}

// Retorna o vizinho que enviou a mensagem para o nó da lista de vizinhos
func RemoveNeighbour(host string, port string, vizinhos []*Vizinho) []*Vizinho {

	new_neighbours := make([]*Vizinho, 0)
	for _, vizinho := range vizinhos {
		if vizinho.HOST != host || vizinho.PORT != port {
			new_neighbours = append(new_neighbours,
				&Vizinho{HOST: vizinho.HOST, PORT: vizinho.PORT})
		}
	}

	return new_neighbours
}

func AdicionaMensagemDFS(noh *No, received_from string, origem_msg string) *DfsMessage {

	temp := &DfsMessage{
		Message:       origem_msg,
		Received_from: received_from,
		Active_child:  "",
		Pending_child: make([]*Vizinho, 0),
	}

	temp.Pending_child = append(temp.Pending_child, noh.Vizinhos...)

	noh.Dfs_messages = append(noh.Dfs_messages, temp)
	return temp
}
