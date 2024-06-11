package main

import (
	"UP2P/client"
	"UP2P/node"
	"UP2P/server"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var SEARCH_FL = client.SearchFlooding
var SEARCH_RW = client.SearchRandomWalk
var SEARCH_DP = client.PrepareSearchInDepth

// Lista de chaves usada para fazer os testes automáticos
var list = [...]string{"Chris_Evans",
	"Kate_McKinnon",
	"Idris_Elba",
	"Julianne_Moore",
	"Ryan_Gosling",
	"Tilda_Swinton",
	"Emma_Watson",
	"Samuel_Jackson",
	"Joaquin_Phoenix",
	"Julia_Roberts",
	"Sandra_Bullock",
	"Brad_Pitt",
	"Emma_Stone",
	"Robert_Downey"}

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

			already_exists := false

			for _, vizinho_na_tabela := range no.Vizinhos {

				if vizinho_na_tabela.HOST == vizinho.HOST &&
					vizinho_na_tabela.PORT == vizinho.PORT {

					already_exists = true

				}

			}

			if !already_exists {
				node.AddNeighbour(no, vizinho.HOST, vizinho.PORT)
			}
		}
	}
}

func listarVizinhos(no *node.No) {
	fmt.Printf("Há %d vizinhos na tabela:\n", len(no.Vizinhos))

	for i, vizinho := range no.Vizinhos {
		fmt.Printf("\t[%d] %s %s\n", i, vizinho.HOST, vizinho.PORT)
	}

}

func calculaDesvioPadrao(array []int, mean float64) float64 {

	square_sum := 0.0

	for _, number := range array {
		square_sum += math.Pow(float64(number)-mean, 2)
	}

	variancia := square_sum / float64(len(array)-1)

	return math.Sqrt(variancia)

}

func showStatistics(no *node.No) {

	fmt.Printf("Estatisticas\n")

	total_flooding := 0
	total_random_walk := 0
	total_search_in_depth := 0

	total_flooding_initialized := 0
	total_random_walk_initialized := 0
	total_search_in_depth_initialized := 0

	total_hops_flooding := 0
	total_hops_random_walk := 0
	total_hops_search_in_depth := 0

	array_flooding := make([]int, 0)
	array_random_walk := make([]int, 0)
	array_search_in_depth := make([]int, 0)

	mean_flooding := 0.0
	mean_random_walk := 0.0
	mean_search_in_depth := 0.0

	standard_deviation_flooding := 0.0
	standard_deviation_random_walk := 0.0
	standard_deviation_search_in_depth := 0.0

	//<ORIGIN> <SEQNO> <TTL> SEARCH <MODE> <LAST_HOP_PORT> <KEY> <HOP_COUNT>

	//<ORIGIN> <SEQNO> <TTL> VAL <MODE> <KEY> <VALUE> <HOP_COUNT>

	for _, mensagem := range no.Received_messages {
		arr := strings.Split(mensagem, " ")

		ACTION := arr[3]
		MODE := arr[4]

		if ACTION == "SEARCH" {
			switch MODE {
			case "FL":
				total_flooding++
			case "RW":
				total_random_walk++
			case "BP":
				total_search_in_depth++
			}
		}

		if ACTION == "VAL" {
			HOPS, _ := strconv.Atoi(arr[7])
			switch MODE {
			case "FL":
				total_flooding_initialized++
				total_hops_flooding += HOPS
				array_flooding = append(array_flooding, HOPS)
			case "RW":
				total_random_walk_initialized++
				total_hops_random_walk += HOPS
				array_random_walk = append(array_random_walk, HOPS)
			case "BP":
				total_search_in_depth_initialized++
				total_hops_search_in_depth += HOPS
				array_search_in_depth = append(array_search_in_depth, HOPS)
			}
		}

	}

	mean_flooding = float64(total_hops_flooding) / float64(total_flooding_initialized)
	mean_random_walk = float64(total_hops_random_walk) / float64(total_random_walk_initialized)
	mean_search_in_depth = float64(total_hops_search_in_depth) / float64(total_search_in_depth_initialized)

	standard_deviation_flooding = calculaDesvioPadrao(array_flooding, mean_flooding)
	standard_deviation_random_walk = calculaDesvioPadrao(array_random_walk, mean_random_walk)
	standard_deviation_search_in_depth = calculaDesvioPadrao(array_search_in_depth, mean_search_in_depth)

	fmt.Printf("\tTotal de mensagens de flooding vistas: %d\n", total_flooding)
	fmt.Printf("\tTotal de mensagens de random walk vistas: %d\n", total_random_walk)
	fmt.Printf("\tTotal de mensagens de busca em profundidade vistas: %d\n", total_search_in_depth)

	fmt.Printf("\tMedia e desvio padrao de saltos ate encontrar destino por flooding: %f %f\n", mean_flooding,
		standard_deviation_flooding)

	fmt.Printf("\tMedia e desvio padrao de saltos ate encontrar destino por random walk: %f %f\n", mean_random_walk,
		standard_deviation_random_walk)

	fmt.Printf("\tMedia e desvio padrao de saltos ate encontrar destino por busca em profundidade: %f %f\n", mean_search_in_depth,
		standard_deviation_search_in_depth)

}

func testesAutomatizados(no *node.No) {

	for i := 1; i <= 3; i++ {
		switch i {
		case 1:
			for _, nome := range list {
				client.FindKey(no, SEARCH_FL, nome)
			}
		case 2:
			for _, nome := range list {
				client.FindKey(no, SEARCH_RW, nome)
			}
		case 3:
			for _, nome := range list {
				client.FindKey(no, SEARCH_DP, nome)
			}
		}
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
			client.FindKey(no, SEARCH_FL, "NONE")
		case 3:
			client.FindKey(no, SEARCH_RW, "NONE")
		case 4:
			client.FindKey(no, SEARCH_DP, "NONE")
		case 5:
			showStatistics(no)
		case 6:
			node.ChangeTTL(no)
		case 7: //Opção secreta para realizar testes automatizados
			testesAutomatizados(no)
		case 9:
			client.Bye(no)
			return
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
		arr := strings.Split(string(data), "\n")
		for _, pair := range arr {
			node.AddKey(pair, noh)
		}
	}

	exibeMenu(noh)
}
