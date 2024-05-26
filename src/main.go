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

var TTL = 100

func imprimeEstadoNo(noh *node.No, mensagem string) {
	fmt.Printf("\n\n\n\n")
	fmt.Printf("////////////////////////// Estado do nó - %s ///////////////////////////////\n", mensagem)

	fmt.Println("Host:", noh.HOST)
	fmt.Println("Port:", noh.PORT)
	fmt.Println("Pares chave-valor: ", noh.Pares_chave_valor)
	fmt.Println("Vizinhos: ", noh.Vizinhos)
	fmt.Println("Mensagens recebidas: ", noh.Mensagens_recebidas)
	fmt.Println("Número de Sequência: ", noh.seqNum)

	fmt.Printf("\n\n\n")
}

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

func GetHOST(address string) string {
	return strings.Split(address, ":")[0]
}

func GetPORT(address string) string {
	return strings.Split(address, ":")[1]
}

func comunicarVizinhos(vizinhos []*node.Vizinho, no *node.No) []string {

	var vizinhosAtivos = make([]string, 0)

	for _, vizinho := range vizinhos {
		time.Sleep(1000 * time.Millisecond)
		status := client.Hello(vizinho.HOST, vizinho.PORT, no)

		if status {
			fmt.Println("Vizinho " + vizinho.HOST + ":" + vizinho.PORT + " sendo adicionado à tabela.\n")
			node.AddNeighbour(no, vizinho.HOST, vizinho.PORT)
		}
	}

	return vizinhosAtivos
}

func exibeMenu(no *node.No) {

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
			client.ShowNeighboursToChoose(no)
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

	if !check_args {
		os.Exit(1)
	}

	address, arqVizinhos, arqParesChaveValor := extraiArgs(args, nArgs)

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

	fmt.Println(paresChaveValor)

	HOST := GetHOST(address)

	PORT := GetPORT(address)

	no := node.NewNo(HOST,
		PORT)

	go server.InitServer(PORT)

	time.Sleep(5000 * time.Millisecond)

	var nosVizinhos []string

	// Envia HELLO para confirmar a existência do vizinho
	if nArgs > 2 {
		nosVizinhos = comunicarVizinhos(string(vizinhos))
	}

	fmt.Println(nosVizinhos)

	exibeMenu(no)
}
