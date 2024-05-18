package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

var TTL = 100

func criaSocketTCP(endereco, porta string) {
	// Cria um socket TCP
	ln, err := net.Listen("tcp4", endereco+":"+porta)
	fmt.Println("Escutando em", endereco+":"+porta)
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
		fmt.Printf("Argumentos: %s\n", args)
		return lenArgs, true
	}

	fmt.Println("Número de argumentos inválido........ >:\n" +
		"Formato: <endereco>:<porta> [vizinhos.txt [lista_chave_valor.txt]]")
	return lenArgs, false
}

func getEnderecoPorta(url string) (string, string) {

	enderecoCompleto := strings.Split(url, ":")
	endereco, porta := enderecoCompleto[0], enderecoCompleto[1]

	return endereco, porta
}

func comunicarVizinhos(vizinhos string) []no {
	vizinhosArr := strings.Split(vizinhos, " ")
	nosVizinhos := make([]no, len(vizinhosArr))

	// range retorna indice, valor. Com "_" estou ignorando o índice no caso abaixo
	for _, vizinho := range vizinhosArr {
		fmt.Println("Tentando adicionar vizinho " + vizinho + ".")

		resposta, err := http.Get(vizinho + "/hello")
		if err != nil {
			fmt.Println("Erro ao enviar a requisição:", err)
			continue
		} else if resposta.StatusCode == 200 {
			var node = new(no)
			nosVizinhos = append(nosVizinhos, *node)
		}
		defer resposta.Body.Close()
	}

	return nosVizinhos
}

func exibeMenu() int {
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
		return -1
	}
	fmt.Println("Número lido:", numero)

	return numero
}

func main() {
	args := os.Args
	nRet, check_args := verificaArgs(args)

	// Não precisa mais por causa do exit
	if !check_args {
		os.Exit(1)
	}

	// Cria socket TCP4 com endereço e porta fornecidos
	endereco, porta := getEnderecoPorta(args[1])
	criaSocketTCP(endereco, porta)

	// Envia HELLO para confirmar a existência do vizinho
	if nRet > 2 {
		nosVizinhos := comunicarVizinhos(args[2])
	}

	comando := exibeMenu()
}
