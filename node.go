package main

import (
	"fmt"
	"strings"
)

type no struct {
	pares_chave_valor map[string]string //Par nome e número associado
	vizinhos          map[string]string //Par endereço e porta
}

func newNo(_pares []string, _vizinhos []string) *no {
	pares := make(map[string]string)

	for i := 0; i < len(_pares)-1; i++ {
		pares[strings.Split(_pares[i], " ")[0]] = strings.Split(_pares[i], " ")[1]
	}

	vizinhos := make(map[string]string)

	for i := 0; i < len(_vizinhos)-1; i++ {
		vizinhos[strings.Split(_vizinhos[i], ":")[0]] = strings.Split(_vizinhos[i], ":")[1]
	}

	fmt.Println(pares)
	fmt.Println(vizinhos)

	return &no{pares_chave_valor: pares, vizinhos: vizinhos}

}
