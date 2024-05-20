package main

import (
	"UP2P/client"
	"UP2P/server"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
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

func extraiArgs(args []string, nArgs int) (string, string, string) {

	if nArgs == 4 {
		return args[1], args[2], args[3]
	}

	if nArgs == 3 {
		return args[2], args[2], ""
	}

	if nArgs == 2 {
		return args[2], "", ""
	}

	return "", "", ""

}

func getEnderecoPorta(url string) (string, string) {

	enderecoCompleto := strings.Split(url, ":")
	endereco, porta := enderecoCompleto[0], enderecoCompleto[1]

	return endereco, porta
}

func comunicarVizinhos(vizinhos string) []string {

	vizinhosArr := strings.Split(vizinhos, "\n")

	var vizinhosAtivos = make([]string, 0)

	for _, vizinho := range vizinhosArr {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("------------------------------------")
		fmt.Printf("Fazendo consulta ao vizinho %s\n", vizinho)
		status := client.Hello(vizinho)

		if status {
			fmt.Println("Vizinho " + vizinho + " sendo adicionado à tabela.\n")
			vizinhosAtivos = append(vizinhosAtivos, vizinho)
		}
	}

	return vizinhosAtivos
}

func exibeMenu(no *no) {

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
			client.ShowNeighboursToChoose(no.vizinhos, client.Hello)
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
	exibeMenu(no)
}

func main() {

	args := os.Args
	nArgs, check_args := verificaArgs(args)

	host, arqVizinhos, arqParesChaveValor := extraiArgs(args, nArgs)

	// Não precisa mais por causa do exit
	if !check_args {
		os.Exit(1)
	}

	var vizinhos string

	if nArgs >= 3 {
		data, status := lerArquivo(arqVizinhos)
		vizinhos = string(data)
		if !status {
			panic("ERRO AO TENTAR ABRIR O ARQUIVO! O PROGRAMA ESTÁ SENDO ENCERRADO...")
		}
	}

	var paresChaveValor string

	if nArgs >= 4 {
		data, status := lerArquivo(arqParesChaveValor)
		paresChaveValor = string(data)
		if !status {
			panic("ERRO AO TENTAR ABRIR O ARQUIVO! O PROGRAMA ESTÁ SENDO ENCERRADO...")
		}

	}

	PORT := ":" + strings.Split(host, ":")[1]

	go server.InitServer(PORT)

	time.Sleep(5000 * time.Millisecond)

	var nosVizinhos []string

	// Envia HELLO para confirmar a existência do vizinho
	if nArgs > 2 {
		nosVizinhos = comunicarVizinhos(string(vizinhos))
	}

	fmt.Println(nosVizinhos)

	no := newNo(host, strings.Split(string(paresChaveValor), "\n"), nosVizinhos)

	exibeMenu(no)
}
