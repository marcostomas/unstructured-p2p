package node

import "strings"

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
}

func NewNo(_HOST string,
	_PORT string) *No {

	return &No{
		HOST:              _HOST,
		PORT:              _PORT,
		NoSeq:             1,
		Pares_chave_valor: map[string]string{},
		Vizinhos:          make([]*Vizinho, 0)}
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
