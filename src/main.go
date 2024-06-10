package main

import (
	"UP2P/client"
	"UP2P/node"
	"UP2P/server"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var SEARCH_FL = client.SearchFlooding
var SEARCH_RW = client.SearchRandomWalk
var SEARCH_DP = client.PrepareSearchInDepth

func lerArquivo(nomeArquivo string) ([]byte, bool) {
	// Abre o arquivo
	arquivo, err := os.Open(nomeArquivo)
	if err != nil {
		return nil, false
	}
	defer arquivo.Close()

	// Lê o conteúdo do arquivo
	conteudo, err := ioutil.ReadAll(arquivo)
	if err != nil {
		return nil, false
	}

	return conteudo, true
}

func verificaArgs(args []string) (int, bool) {
	lenArgs := len(args)

	if lenArgs > 0 && lenArgs < 5 {
		return lenArgs, true
	}

	return lenArgs, false
}

func extraiArgs(args []string, nArgs int) (string, string, string) {

	if nArgs == 4 {
		return args[1], args[2], args[3]
	}

	if nArgs == 3 {
		return args[1], args[2], ""
	}

	if nArgs == 2 {
		return args[1], "", ""
	}

	return "", "", ""

}

func comunicarVizinhos(vizinhos []*node.Vizinho, no *node.No) {

	for _, vizinho := range vizinhos {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("Tentando adicionar " + vizinho.HOST + ":" + vizinho.PORT + ".\n")
		status := client.Hello(vizinho.HOST, vizinho.PORT, no)

		if status {
			node.AddNeighbour(no, vizinho.HOST, vizinho.PORT)
		}
	}
}

func listarVizinhos(no *node.No) {
	fmt.Printf("Há %d vizinhos na tabela:\n", len(no.Vizinhos))

	for i, vizinho := range no.Vizinhos {
		fmt.Printf("\t[%d] %s %s\n", i, vizinho.HOST, vizinho.PORT)
	}

}

func exibeMenu(no *node.No) {
	fmt.Println("")
	fmt.Println("Escolha o comando")
	fmt.Println("[0] Listar vizinhos")
	fmt.Println("[1] Hello")
	fmt.Println("[2] SEARCH (flooding)")
	fmt.Println("[3] SEARCH (random walk)")
	fmt.Println("[4] SEARCH (busca em profundidade)")
	fmt.Println("[5] Estatísticas")
	fmt.Println("[6] Alterar valor padrão de TTL")
	fmt.Println("[9] Sair")

	var numero int
	_, err := fmt.Scanln(&numero)

	if err != nil {
		fmt.Println("Erro ao ler o número", err)
	} else {

		switch numero {
		case 0:
			listarVizinhos(no)
		case 1:
			client.ShowNeighbours(no)
		case 2:
			client.FindKey(no, SEARCH_FL)
		case 3:
			client.FindKey(no, SEARCH_RW)
		case 4:
			client.FindKey(no, SEARCH_DP)
		case 5:
			fmt.Println("Estatísticas")
		case 6:
			node.ChangeTTL(no)
		case 9:

		}
	}
	exibeMenu(no)
}

func main() {

	args := os.Args
	nArgs, check_args := verificaArgs(args)

	if !check_args {
		os.Exit(1)
	}

	address, arqVizinhos, arqParesChaveValor := extraiArgs(args, nArgs)

	HOST := strings.Split(address, ":")[0]
	PORT := strings.Split(address, ":")[1]

	noh := node.NewNo(HOST, PORT)
	go server.InitServer(noh)

	time.Sleep(5000 * time.Millisecond)
	// Envia HELLO para confirmar a existência do vizinho
	if nArgs > 2 {
		data, status := lerArquivo(arqVizinhos)
		if !status {
			panic("ERRO AO TENTAR ABRIR O ARQUIVO! O PROGRAMA ESTÁ SENDO ENCERRADO...")
		}
		listaVizinhos := node.GenerateNeighboursList(data)
		comunicarVizinhos(listaVizinhos, noh)
	}

	if nArgs == 4 {
		data, status := lerArquivo(arqParesChaveValor)
		if !status {
			panic("ERRO AO TENTAR ABRIR O ARQUIVO! O PROGRAMA ESTÁ SENDO ENCERRADO...")
		}
		fmt.Printf("Tentando adicionar pares chave valor na tabela...\n")
		arr := strings.Split(string(data), "\n")
		for _, pair := range arr {
			node.AddKey(pair, noh)
		}
	}

	exibeMenu(noh)
}
