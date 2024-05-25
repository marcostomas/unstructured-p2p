# Algumas keywords "específicas" de Go

**go**: a declaração de "go" inicia uma execução de uma chamada de função como uma thread concorrente independente de controle (aka goroutine).

**chan**: nos ajuda a comunicar entre goroutines. esta palavra nos permite definir um canal, o que nos dá uma forma econômica de enviar e receber valores entre goroutines. com canais podemos sincronizar goroutines, ao passar valores.

**select**: é usado com goroutines para lidar com múltiplas operações de canais. esta declaração permite um programa esperar em múltiplas operações de comunicaçã0 (canais), seguindo a execução do programa com o canal que primeiro estiver pronto.

**defer**: garante que uma chamada de função possa ser executada mais tarde, durante a execução do programa.

```go
func main(){
    defer fmt.Println("world")
    fmt.Println("hello")
}
// hello
// world
```

**range**: usado para iterar sobre elementos de diferentes tipos de slices.

**fallthrough**: indica, dentro de um switch, que a linha de execução deve passar para o _case_ abaixo.
