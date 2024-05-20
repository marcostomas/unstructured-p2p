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

func verificaArgs(args []string) (int, bool) {
	lenArgs := len(args)

	if lenArgs > 0 || lenArgs < 5 {
		return lenArgs, true
	}

	fmt.Println("Número de argumentos inválido........ >:\n" +
		"Formato: <endereco>:<porta> [vizinhos.txt [lista_chave_valor.txt]]")
	return -1, false
}

func main() {
	args := os.Args
	var endereco, porta string
	nRet, check_args := verificaArgs(args)

	//Só pra parar de dar warning de unused variable
	fmt.Println(nRet)

	// Não precisa mais por causa do exit
	if !check_args {
		os.Exit(1)
	}

	enderecoCompleto := strings.Split(args[1], ":")
	endereco, porta = enderecoCompleto[0], enderecoCompleto[1]

	criaSocketTCP(endereco, porta)

	init_server()

	init_client()

}
