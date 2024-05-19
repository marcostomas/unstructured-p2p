package main

import (
	"UP2P/server"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var TTL = 100
var node *no

type no struct {
	pares_chave_valor   map[string]string
	vizinhos            []string
	mensagens_recebidas []string
	seqNum              int
}

func imprimeEstadoNo(noh *no, mensagem string) {
	fmt.Printf("\n\n\n\n\n\n")
	fmt.Printf("//////////////////////////Estado do nó - %s///////////////////////////////", mensagem)

	fmt.Println("Pares chave-valor: ", noh.pares_chave_valor)
	fmt.Println("Vizinhos: ", noh.vizinhos)
	fmt.Println("Mensagens recebidas: ", noh.mensagens_recebidas)

	fmt.Printf("\n\n\n")
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

	if lenArgs > 1 && lenArgs < 5 {
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

func inicializaNode() *no {
	node = new(no)
	node.pares_chave_valor = make(map[string]string)
	node.vizinhos = make([]string, 0)
	node.mensagens_recebidas = make([]string, 0)
	node.seqNum = 1 // Primeiro envia a mensagem, depois incrementa o seqNum
	return node
}

func comunicarVizinhos(vizinhos string, noh *no) {
	vizinhosArr := strings.Split(vizinhos, " ")

	// range retorna indice, valor. Com "_" estou ignorando o índice no caso abaixo
	for _, vizinho := range vizinhosArr {
		fmt.Println("Tentando adicionar vizinho " + vizinho + ".")

		resposta, err := http.Get("http://" + vizinho + "/hello")
		if err != nil {
			fmt.Println("Erro ao enviar a requisição:", err)
			continue
		} else if resposta.StatusCode == 200 {
			fmt.Println("Envio feito com sucesso: ")

			noh.vizinhos = append(noh.vizinhos, vizinho)
		}
		defer resposta.Body.Close()
	}
}

func alterarTTL() {
	fmt.Println("Digite o novo valor de TTL:")
	_, err := fmt.Scanln(&TTL)
	if err != nil {
		fmt.Println("Erro ao ler o número", err)
		alterarTTL()
	}
	fmt.Println("TTL alterado para:", TTL)
}

func sair(noh *no) {
	fmt.Println("Saindo...")
	// Envia mensagem de saída para todos os vizinhos
	for _, vizinho := range noh.vizinhos {
		client.bye(vizinho)
	}

	os.Exit(0)
}

func hello() {
	/*
		for _, vizinho := range node.nosVizinhos{
			client.Hello(getEnderecoPorta(vizinho))
		}
	*/
}

func adicionaChaveDoNo(pares string, noh *no) {
	paresArr := strings.Split(pares, "\n")

	for _, par := range paresArr {
		parArr := strings.Split(par, " ")
		noh.pares_chave_valor[parArr[0]] = parArr[1]
	}
}

func exibeMenu(noh *no) {
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
		exibeMenu(noh)
	}
	fmt.Println("Você escolheu ", numero)

	switch numero {
	case 0:
		fmt.Println("Listar vizinhos")
	case 1:
		hello()
	case 2:
		fmt.Println("SEARCH (flooding)")
	case 3:
		fmt.Println("SEARCH (random walk)")
	case 4:
		fmt.Println("SEARCH (busca em profundidade)")
	case 5:
		fmt.Println("Estatísticas")
	case 6:
		alterarTTL()
	case 9:
		sair(noh)
	}

	exibeMenu(noh)
}

func main() {
	args := os.Args
	nArgs, check_args := verificaArgs(args)

	// Não precisa mais por causa do exit
	if !check_args {
		os.Exit(1)
	}

	// Cria socket TCP4 com endereço e porta fornecidos
	endereco, porta := getEnderecoPorta(args[1])
	go server.InitServer(endereco, porta)

	noh := inicializaNode()

	imprimeEstadoNo(noh, "1")
	// Envia HELLO para confirmar a existência do vizinho
	if nArgs > 2 {
		comunicarVizinhos(string(lerArquivo(args[2])), noh)
	}

	imprimeEstadoNo(noh, "2")

	if nArgs == 4 {
		adicionaChaveDoNo(string(lerArquivo(args[3])), noh)
	}

	imprimeEstadoNo(noh, "3")

	exibeMenu()
}
