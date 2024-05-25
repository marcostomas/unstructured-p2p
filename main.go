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
var node *no

func imprimeEstadoNo(noh *no, mensagem string) {
	fmt.Printf("\n\n\n\n")
	fmt.Printf("////////////////////////// Estado do nó - %s ///////////////////////////////\n", mensagem)

	fmt.Println("Host:", noh.HOST)
	fmt.Println("Port:", noh.PORT)
	fmt.Println("Pares chave-valor: ", noh.pares_chave_valor)
	fmt.Println("Vizinhos: ", noh.vizinhos)
	fmt.Println("Mensagens recebidas: ", noh.mensagens_recebidas)
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

	if lenArgs > 0 && lenArgs < 5 {
		fmt.Println("Validação OK!")
		fmt.Printf("Argumentos: %s\n", args)
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
		return args[1], args[2], ""
	}

	if nArgs == 2 {
		return args[1], "", ""
	}

	return "", "", ""

}

func getEnderecoPorta(url string) (string, string) {
	enderecoCompleto := strings.Split(url, ":")
	endereco, porta := enderecoCompleto[0], enderecoCompleto[1]

	return endereco, porta
}

func comunicarVizinhos(vizinhos string, noh *no) {
	vizinhosArr := strings.Split(vizinhos, "\n")

	// range retorna indice, valor. Com "_" estou ignorando o índice no caso abaixo
	for _, vizinho := range vizinhosArr {
		fmt.Println("Tentando adicionar vizinho " + vizinho + ".")

		ret := client.Hello(noh.HOST, noh.PORT, string(noh.seqNum), string(TTL), "Hello", vizinho, noh.PORT)
		if ret == true {
			fmt.Println("Envio feito com sucesso: ")
			noh.vizinhos = append(noh.vizinhos, vizinho)
			continue
		} else {
			fmt.Println("Erro ao enviar a requisição")
			continue
		}
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
	// for _, vizinho := range noh.vizinhos {
	// 	client.bye(vizinho)
	// }

	os.Exit(0)
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
		fmt.Println("Listar Vizinhos")
		// client.ShowNeighboursToChoose(noh.vizinhos)
	case 1:
		client.Hello(noh.HOST, noh.PORT, string(noh.seqNum), string(TTL), "Hello", noh.HOST, noh.PORT)
	case 2:
		client.SearchFlooding("")
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
	host, arqVizinhos, arqParesChaveValor := extraiArgs(args, nArgs)

	noh := inicializaNode()
	noh.HOST, noh.PORT = getEnderecoPorta(host)
	go server.InitServer(noh.HOST, noh.PORT)

	time.Sleep(5000 * time.Millisecond)

	imprimeEstadoNo(noh, "1")
	// Envia HELLO para confirmar a existência do vizinho
	if nArgs > 2 {
		data, status := lerArquivo(arqVizinhos)
		if !status {
			panic("ERRO AO TENTAR ABRIR O ARQUIVO! O PROGRAMA ESTÁ SENDO ENCERRADO...")
		}
		comunicarVizinhos(string(data), noh)
	}

	imprimeEstadoNo(noh, "2")

	if nArgs == 4 {
		data, status := lerArquivo(arqParesChaveValor)
		if !status {
			panic("ERRO AO TENTAR ABRIR O ARQUIVO! O PROGRAMA ESTÁ SENDO ENCERRADO...")
		}
		adicionaChaveDoNo(string(data), noh)
	}

	imprimeEstadoNo(noh, "3")

	exibeMenu(noh)
}
