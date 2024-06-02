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
}

func NewNo(_HOST string, _PORT string) *No {

	return &No{
		HOST:              _HOST,
		PORT:              _PORT,
		NoSeq:             1,
		Pares_chave_valor: map[string]string{},
		Vizinhos:          make([]*Vizinho, 0),
		Received_messages: make([]string, 0),
		TTL:               1,
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
	fmt.Printf("\n")
	fmt.Printf("/////////////////// Estado do nó - %d ///////////////////\n", count)

	fmt.Println("HOST: ", noh.HOST)
	fmt.Println("PORT: ", noh.PORT)
	fmt.Println("Pares chave-valor: ", noh.Pares_chave_valor)
	fmt.Println("Vizinhos: ", noh.Vizinhos)
	fmt.Println("Mensagens recebidas: ", noh.Received_messages)
	fmt.Println("Número de Sequência: ", noh.NoSeq)

	fmt.Printf("\n")
}

func FindReceivedMessage(message string, NO *No) bool {
	for _, msg := range NO.Received_messages {
		if msg == message {
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
func RemoveNeighbour(host string, port string, no *No) []*Vizinho {

	new_neighbours := make([]*Vizinho, len(no.Vizinhos))
	for _, vizinho := range no.Vizinhos {
		if vizinho.HOST != host && vizinho.PORT != port {
			new_neighbours = append(new_neighbours,
				&Vizinho{HOST: vizinho.HOST, PORT: vizinho.PORT})
		}
	}

	return new_neighbours
}
