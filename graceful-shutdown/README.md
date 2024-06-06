# graceful-shutdown

Validação de maneiras diferentes de fazer um shutdown de uma aplicação em Go.

## Estrutura

### v1

Estrutura menos convencional que encapsula o shutdown em uma goroutine e o listen do servidor no processo principal.

### v2

Estrutura mais comum que encapsula o listen do servidor em uma goroutine e o shutdown no processo principal.

### v3

Estrutura semelhante a v1, porém com um canal pra bloquear o processo principal de ser finalizado.

## Testes

Para testar as diferentes maneiras de fazer um shutdown, execute o seguinte comando:

```bash
$ go run main.go
```

O script sobe todas as versões (v1, v2, v3), uma por vez, realiza uma requisição e testa o graceful shutdown.

## Resultados

- v1: Não funciona, já que ignora a requisição em andamento.
- v2: Funciona bem, já que aguarda a requisição em andamento por até 10 segundos.
- v3: Funciona bem, já que aguarda a requisição em andamento por até 10 segundos. Porém, é fácil de esquecer de tratar o canal de bloqueio.
