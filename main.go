package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func criaSocketTCP(endereco, porta string) {
	// Cria um socket TCP
	ln, err := net.Listen("tcp4", endereco+":"+porta)
	if err != nil {
		fmt.Println("Erro ao criar socket TCP")
		return
	}
	defer ln.Close()
}

func lerArquivo(nomeArquivo string) []byte {
	// Abre o arquivo
	arquivo, err := os.Open(nomeArquivo)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return nil
	}
	defer arquivo.Close()

	// Lê o conteúdo do arquivo
	conteudo, err := ioutil.ReadAll(arquivo)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return nil
	}

	return conteudo
}

func main() {
	// os.Args fornece acesso aos argumentos de linha de comando
	args := os.Args
	var endereco, porta string
	if len(args) < 1 {
		fmt.Println("Informe um endereço e porta para criar o socket TCP")
	} else if len(args) == 2 {
		fmt.Println("Endereço e porta informados: ", args[1])
		enderecoCompleto := strings.Split(args[1], ":")
		endereco, porta = enderecoCompleto[0], enderecoCompleto[1]

	} else if len(args) == 4 {
		fmt.Println(fmt.Sprintf("Endereço e porta informados: %s \n", args[1]))

		vizinhosRaw := lerArquivo(args[2])
		vizinhos := strings.Split(string(vizinhosRaw), "\n")

		fmt.Println("Vizinhos: ", vizinhos)

		paresRaw := lerArquivo(args[3])
		pares := strings.Split(string(paresRaw), "\n")

		fmt.Println("pares:", pares)
	} else {
		fmt.Println("Número de argumentos inválido")
		return
	}

	criaSocketTCP(endereco, porta)
}
