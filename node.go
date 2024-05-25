package main

import "strings"

type no struct {
	pares_chave_valor   map[string]string //Par nome e número associado
	vizinhos            []string
	mensagens_recebidas []string
	seqNum              int
	HOST                string
	PORT                string
}

var node *no

func inicializaNode() *no {
	node = new(no)
	node.pares_chave_valor = make(map[string]string)
	node.vizinhos = make([]string, 0)
	node.mensagens_recebidas = make([]string, 0)
	node.seqNum = 1 // Primeiro envia a mensagem, depois incrementa o seqNum
	return node
}

func adicionaChaveDoNo(pares string, noh *no) {
	paresArr := strings.Split(pares, "\n")

	for _, par := range paresArr {
		parArr := strings.Split(par, " ")
		noh.pares_chave_valor[parArr[0]] = parArr[1]
	}
}
