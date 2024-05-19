- O main será executado para cada nó

# Inicialização

- 1º abrir um socket TCP (IPv4)
- Quando um nó recebe uma requisição no endpoint /hello, ele deve verificar se quem está mandando já está na sua tabela de vizinhos

# Checklist

| Feito? | Tarefa                                            | Responsável | Prazo |
| ------ | ------------------------------------------------- | ----------- | ----- |
| [X]    | Inicialização                                     | Marcos      | 19/05 |
| [X]    | Menu, controle de comando, alteração de TTL, sair | Marcos      | 19/05 |
| [ ]    | Listar vizinhos                                   | Rafael      | 19/05 |
| [ ]    | Hello                                             | Rafael      | 19/05 |
| [ ]    | Search (flooding)                                 |             |       |
| [ ]    | Search (random walk)                              |             |       |
| [ ]    | Search (busca em profundidade)                    |             |       |
| [ ]    | Estatísticas                                      |             |       |
| [ ]    | Relatório                                         |             |       |

# Endpoints do servidor

Para cada busca teremos dois endpoints, um para pedir que o próximo nó verifique em sua tabela interna; e o outro para comunicar uma chave encontrada.

Busca por flooding: /EncontradoFL, /BuscarFL
Busca por random walk: /EncontradoRW, BuscarRW
Busca por profundidade: /EncontradoBP, BuscarBP
Hello
Bye

Total de 8 endpoints

# Campos das mensagens

- ORIGIN: endereço e porta do nó que iniciou a busca
- SEQNO: número de sequência. é um contador por nó que indica quantas mensagens já foram enviadas por aquele nó.
- TTL: time to live. representa o número máximo de saltos que são toleráveis até encontrar a chave; se não for encontrada, o pacote morre
- MODE: o modo de busca. pode ser `FL`, `RW` ou `BP`
- LAST_HOP_PORT: indica a porta do servidor que encaminhou a mensagem.
- KEY: é a chave que deve ser procurada na tabela de cada nó.
- HOP_COUNT: contagem de saltos a serem realizados. HOP_COUNT <= TTL. HOP_COUNT é 1 no primeiro envio
- NEIGH: é um dos vizinhos do nó que está encaminhando a mensagem. formato: endereco:porta

# Endpoints

### /Buscar\*

Método: `GET`
Body: recebe uma mensagem do tipo `<ORIGIN> <SEQNO> <TTL> SEARCH <MODE> <LAST_HOP_PORT> <KEY> <HOP_COUNT>\n`
Retorno: `<MODE>_OK`, ao receber uma mensagem válida.

Imprime na tela: `Encaminhando mensagem "<ORIGIN> <SEQNO> <TTL> <KEY>" para <NEIGH>`. Pode ser implementado no cliente também
Imprime na tela 2: `Envio feito com sucesso: "<NEIGH> <SEQNO> <TTL> <KEY>"`. Este aparece apenas quando `<NEIGH>` confirma o recebimento da mensagem. Pode ser implementado no cliente também.

### /Encontrado\*

Método: `GET`
Body: receby uma mensagem do tipo `<ORIGIN> <SEQNO> <TTL> VAL <MODE> <KEY> <VALUE> <HOP_COUNT>`
Retorna: `TRUE`, ao receber uma mensagem válida.

Imprime na tela: `Valor encontrado! Chave: <KEY> valor: <VALOR>`

- Aqui, `<HOP_COUNT>` é igual ao valor que foi recebido, no seu endpoint /BuscarFL, pelo nó que encontrou a chave.

### /Hello

Método: `GET`
Body: mensagem do tipo `<ORIGIN> <SEQNO> <TTL> HELLO`
Retorno: `TRUE`, se o TTL é igual a `1`

Imprime na tela: `Mensagem recebida: "<ORIGIN> <SEQNO> <TTL> HELLO"`

1: `Adicionando vizinho na tabela: <ORIGIN>`
2: `Vizinho já está na tabela: <ORIGIN>`

- TTL de mensagens enviadas para este endpoint devem ser igual a `1`
- Ao receber uma requisição neste endpoint, e o TTL ser válido (igual a `1`), o nó deve adicionar o remetente da mensagem na sua tabela de vizinhos. Se o remetente não está na tabela, imprime 1; se estiver, imprime 2.

### /Bye

Método: `GET`
Body: mensagem do tipo `<ORIGIN> <SEQNO> <TTL> BYE`
Retorno: `TRUE`

Imprime na tela: `Mensagem recebida: "<ORIGIN> <SEQNO> <TTL> BYE"`

- Ao receber uma requisição neste endpoint, o nó remove o remetente desta requisição da sua tabela de vizinhos
- Depois de remover, imprime na tela: `Removendo vizinho da tabela: <ORIGIN>`
