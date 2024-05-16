package main

import (
	"fmt"
	"net"
	"os"
	"io/ioutil"
)

func criaSocketTCP(){
	// Cria um socket TCP
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao criar socket TCP")
		return
	}
	defer ln.Close()
}

func lerArquivo(nomeArquivo string) {
    // Abre o arquivo
    arquivo, err := os.Open(nomeArquivo)
    if err != nil {
        fmt.Println("Erro ao abrir o arquivo:", err)
        return
    }
    defer arquivo.Close()

    // Lê o conteúdo do arquivo
    conteudo, err := ioutil.ReadAll(arquivo)
    if err != nil {
        fmt.Println("Erro ao ler o arquivo:", err)
        return
    }

    // Imprime o conteúdo do arquivo
    fmt.Println(string(conteudo))
	fmt.Println("")
}

func main(){
	// os.Args fornece acesso aos argumentos de linha de comando
    args := os.Args

    // os.Args[0] é o nome do programa
    // os.Args[1:] são os argumentos que foram passados
    fmt.Println("Nome do programa:", args[0])
    fmt.Println("Argumentos passados:", args[1:])


	if len(args) < 1 {
		fmt.Println("Informe um endereço e porta para criar o socket TCP");
	} else if len(args) == 2 {
		fmt.Println("Endereço e porta informados: ", args[1]);
	} else if len(args) == 4{
		fmt.Println(fmt.Sprintf("Endereço e porta informados: %s \n", args[1]));

		fmt.Println("Arquivo de vizinhos:");
		lerArquivo(args[2]);

		fmt.Println("Arquivo de pares chave-valor:");
		lerArquivo(args[3]);
	} else {
		fmt.Println("Número de argumentos inválido");
		return
	}

	// lerTerminal()
	// criaSocketTCP()
}