package main

import (
	"UP2P/client"
	"UP2P/server"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var TTL = 100

func lerArquivo(nomeArquivo string) ([]byte, bool) {
	// Abre o arquivo
	arquivo, err := os.Open(nomeArquivo)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return nil, false
	}
	defer arquivo.Close()

	// Lê o conteúdo do arquivo
	conteudo, err := ioutil.ReadAll(arquivo)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return nil, false
	}

	return conteudo, true
}

func verificaArgs(args []string) (int, bool) {
	lenArgs := len(args)

	fmt.Println("------------------------------------")
	fmt.Println("Fazendo verificação da quantidade de argumentos...")

	if lenArgs > 0 || lenArgs < 5 {
		fmt.Println("Validação OK!")
		fmt.Printf("Argumentos: %s\n\n", args)
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

func exibeMenu() {

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
			client.SearchFlooding("")
		case 1:
			client.Hello()
		case 2:
			fmt.Println("SEARCH (flooding)")
		case 3:
			fmt.Println("SEARCH (random walk)")
		case 4:
			fmt.Println("SEARCH (busca em profundidade)")
		case 5:
			fmt.Println("Estatísticas")
		case 6:
			fmt.Println("Alterar valor padrão de TTL")
		case 9:
			os.Exit(0)
		}
	}
	exibeMenu()
}

func main() {

	args := os.Args
	nArgs, check_args := verificaArgs(args)

	// Não precisa mais por causa do exit
	if !check_args {
		os.Exit(1)
	}

	fmt.Println("------------------------------------")
	fmt.Println("Tentando extrair os dados de " + args[2] + "...")
	vizinhos, status := lerArquivo(args[2])
	if !status {
		panic("ERRO AO TENTAR ABRIR O ARQUIVO! O PROGRAMA ESTÁ SENDO ENCERRADO...")
	} else {
		fmt.Println("Leitura do arquivo bem sucedida!")
	}

	fmt.Println("------------------------------------")
	fmt.Println("Tentando extrair os dados de " + args[3] + "...")
	paresChaveValor, status := lerArquivo(args[3])
	if !status {
		panic("ERRO AO TENTAR ABRIR O ARQUIVO! O PROGRAMA ESTÁ SENDO ENCERRADO...")
	} else {
		fmt.Println("Leitura do arquivo bem sucedida!")
	}

	fmt.Println("------------------------------------")
	fmt.Println("Vizinhos: ")
	fmt.Println(string(vizinhos) + "\n")

	fmt.Println("------------------------------------")
	fmt.Println("Pares chave-valor: ")
	fmt.Println(string(paresChaveValor) + "\n")

	no := newNo(strings.Split(string(paresChaveValor), ";"),
		strings.Split(string(vizinhos), ";"))

	fmt.Println(no.pares_chave_valor)
	fmt.Println(no.vizinhos)

	// Envia HELLO para confirmar a existência do vizinho
	if nArgs > 2 {
		nosVizinhos := comunicarVizinhos(args[2])

		fmt.Println(nosVizinhos)

	}

	go server.InitServer()

	exibeMenu()
}
