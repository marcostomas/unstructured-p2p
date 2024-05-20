package main

import (
	"strings"
)

type no struct {
	host              string
	seqNo             int
	pares_chave_valor map[string]string //Par nome e número associado
	vizinhos          []string          //Endereço:porta
}

func newNo(_host string, _pares []string, _vizinhos []string) *no {
	_seqNo := 1

	pares := make(map[string]string)

	for i := 0; i < len(_pares); i++ {
		pares[strings.Split(_pares[i], " ")[0]] = strings.Split(_pares[i], " ")[1]
	}

	return &no{host: _host, seqNo: _seqNo, pares_chave_valor: pares, vizinhos: _vizinhos}

}
