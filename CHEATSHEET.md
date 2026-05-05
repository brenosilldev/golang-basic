# Cheatsheet Completo de Go (Golang)

Guia de referência rápida e explicativa da linguagem Go. Use o sumário para navegar até o tópico desejado.

## Sumário

1. [Estrutura básica de um programa](#1-estrutura-básica-de-um-programa)
2. [Comandos da CLI `go`](#2-comandos-da-cli-go)
3. [Módulos e pacotes](#3-módulos-e-pacotes)
4. [Variáveis, constantes e tipos](#4-variáveis-constantes-e-tipos)
5. [Operadores](#5-operadores)
6. [Strings, runes e bytes](#6-strings-runes-e-bytes)
7. [Controle de fluxo](#7-controle-de-fluxo)
8. [Funções](#8-funções)
9. [Arrays, slices e maps](#9-arrays-slices-e-maps)
10. [Structs](#10-structs)
11. [Métodos e interfaces](#11-métodos-e-interfaces)
12. [Ponteiros](#12-ponteiros)
13. [Tratamento de erros](#13-tratamento-de-erros)
14. [Defer, panic e recover](#14-defer-panic-e-recover)
15. [Goroutines e canais (concorrência)](#15-goroutines-e-canais-concorrência)
16. [Generics](#16-generics)
17. [Testes e benchmarks](#17-testes-e-benchmarks)
18. [Pacotes da biblioteca padrão mais usados](#18-pacotes-da-biblioteca-padrão-mais-usados)
19. [Convenções e boas práticas](#19-convenções-e-boas-práticas)

---

## 1. Estrutura básica de um programa

Todo programa Go executável precisa de:

- O pacote `main`.
- Uma função `main()` sem parâmetros e sem retorno.
- Imports declarados no topo do arquivo.

```go
package main

import "fmt"

func main() {
    fmt.Println("Olá, Mundo!")
}
```

Pontos importantes:

- Cada arquivo `.go` começa com `package <nome>`. Apenas `package main` produz binário.
- Imports não usados e variáveis não usadas geram **erro de compilação** (não warning).
- Não há ponto e vírgula no fim das linhas (o compilador insere automaticamente).
- A chave `{` precisa estar na **mesma linha** da declaração (`func`, `if`, etc.).

---

## 2. Comandos da CLI `go`

| Comando | Para que serve |
|---|---|
| `go run arquivo.go` | Compila e executa em um passo, sem gerar binário permanente. |
| `go build` | Compila e gera binário no diretório atual. |
| `go install` | Compila e instala o binário em `$GOBIN`. |
| `go mod init <path>` | Inicializa um módulo (cria `go.mod`). |
| `go mod tidy` | Adiciona/remove dependências de acordo com o código. |
| `go get <pacote>` | Baixa e adiciona dependência ao módulo. |
| `go test ./...` | Roda todos os testes do módulo. |
| `go test -v -run NomeDoTeste` | Roda um teste específico em modo verboso. |
| `go test -bench=. -benchmem` | Executa benchmarks com info de memória. |
| `go fmt ./...` | Formata todos os arquivos no padrão oficial. |
| `go vet ./...` | Análise estática para encontrar erros comuns. |
| `go doc fmt.Println` | Mostra documentação de um símbolo. |
| `go env` | Mostra variáveis de ambiente do Go. |

---

## 3. Módulos e pacotes

**Módulo** é o conjunto de pacotes versionado por um `go.mod`.
**Pacote** é a pasta com arquivos `.go` que compartilham o mesmo `package`.

`go.mod` típico:

```go
module github.com/usuario/projeto

go 1.22

require (
    github.com/pkg/errors v0.9.1
)
```

Importando pacotes locais:

```go
import (
    "fmt"                                   // padrão
    "github.com/usuario/projeto/modulo1"    // pacote local do módulo
    m2 "github.com/usuario/projeto/modulo2" // alias
    _ "image/png"                           // só side-effects (init)
    . "math"                                // expõe sem prefixo (evite!)
)
```

**Visibilidade:** identificadores que começam com **letra maiúscula** são exportados (públicos); minúscula são privados ao pacote.

```go
func Soma(a, b int) int { return a + b } // exportada
func soma(a, b int) int { return a + b } // privada
```

---

## 4. Variáveis, constantes e tipos

Declaração de variáveis:

```go
var nome string = "Breno"   // forma completa
var idade = 30              // tipo inferido
var ativo bool              // valor zero (false)

nome := "Breno"             // forma curta (apenas dentro de função)

var (                       // bloco de declarações
    x int    = 10
    y string = "olá"
)
```

**Valores zero** (default): `0` para numéricos, `""` para string, `false` para bool, `nil` para ponteiros, slices, maps, channels, funcs e interfaces.

Constantes:

```go
const Pi = 3.14159
const (
    A = iota // 0
    B        // 1
    C        // 2
)
```

Tipos primitivos:

| Categoria | Tipos |
|---|---|
| Inteiros com sinal | `int`, `int8`, `int16`, `int32`, `int64` |
| Inteiros sem sinal | `uint`, `uint8` (=`byte`), `uint16`, `uint32`, `uint64`, `uintptr` |
| Ponto flutuante | `float32`, `float64` |
| Complexos | `complex64`, `complex128` |
| Texto | `string`, `rune` (=`int32`) |
| Booleano | `bool` |

Conversão de tipos é **explícita**:

```go
i := 42
f := float64(i)
s := strconv.Itoa(i) // int -> string
n, err := strconv.Atoi("42")
```

---

## 5. Operadores

- **Aritméticos:** `+ - * / %` e `++ --` (apenas como statement, ex.: `i++`).
- **Comparação:** `== != < <= > >=`.
- **Lógicos:** `&& || !`.
- **Bit a bit:** `& | ^ &^ << >>`.
- **Atribuição:** `= += -= *= /= %= &= |= ^= <<= >>=`.
- **Endereço/desreferência:** `&x`, `*p`.
- **Canal:** `<-`.

`&^` é o "AND NOT" (bit clear), exclusivo de Go.

---

## 6. Strings, runes e bytes

- `string` é uma sequência **imutável** de bytes (UTF-8).
- `byte` é alias de `uint8` (um byte).
- `rune` é alias de `int32` (um code point Unicode).

```go
s := "olá"
fmt.Println(len(s))        // 4 (bytes, não caracteres!)
for i, r := range s {      // i = byte index, r = rune
    fmt.Printf("%d %c\n", i, r)
}

b := []byte(s)             // converte para slice de bytes
r := []rune(s)             // converte para slice de runes
```

Operações comuns (pacote `strings`):

```go
strings.Contains("foobar", "oo")   // true
strings.HasPrefix("hello", "he")   // true
strings.Split("a,b,c", ",")        // ["a","b","c"]
strings.Join([]string{"a","b"},"-")// "a-b"
strings.ToUpper("go")              // "GO"
strings.ReplaceAll("aaa", "a","b") // "bbb"
strings.TrimSpace("  hi  ")        // "hi"
```

Formatação (`fmt.Printf`):

| Verbo | Significado |
|---|---|
| `%v` | valor padrão |
| `%+v` | struct com nomes dos campos |
| `%#v` | sintaxe Go |
| `%T` | tipo |
| `%d` | inteiro decimal |
| `%b %o %x %X` | binário, octal, hex |
| `%f %e %g` | float |
| `%s` | string |
| `%q` | string com aspas |
| `%c` | rune como caractere |
| `%p` | ponteiro |

---

## 7. Controle de fluxo

**`if` / `else`** (com inicializador opcional):

```go
if n, err := strconv.Atoi("42"); err == nil {
    fmt.Println(n)
} else {
    fmt.Println(err)
}
```

**`for`** é o único laço — substitui `while` e `do-while`:

```go
for i := 0; i < 10; i++ { ... }   // clássico
for cond { ... }                   // como while
for { ... }                        // infinito
for i, v := range slice { ... }    // iterar slice/array/map/string/channel
```

`break`, `continue` e `goto label` funcionam como esperado. Pode-se usar **labels**:

```go
Outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 2 { break Outer }
    }
}
```

**`switch`** — não precisa de `break`; cada `case` quebra automaticamente. Use `fallthrough` para continuar.

```go
switch x := dia; x {
case "seg", "ter":
    fmt.Println("início")
case "sex":
    fmt.Println("sextou")
default:
    fmt.Println("outro")
}

switch { // sem expressão = if/else if encadeado
case n < 0:
case n == 0:
case n > 0:
}

switch v := i.(type) { // type switch
case int: ...
case string: ...
}
```

---

## 8. Funções

```go
func soma(a, b int) int { return a + b }

// múltiplos retornos
func divmod(a, b int) (int, int) { return a / b, a % b }

// retornos nomeados (já inicializados com valor zero)
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // "naked return"
}

// variádica
func soma(ns ...int) int {
    total := 0
    for _, n := range ns { total += n }
    return total
}
soma(1, 2, 3)
soma(nums...) // expandir slice

// função como valor / closure
add := func(a, b int) int { return a + b }
contador := func() func() int {
    n := 0
    return func() int { n++; return n }
}()
```

Argumentos são passados **por valor** (cópia). Para mutar o original passe um ponteiro ou use slices/maps (cuja cabeça é cópia, mas referencia o mesmo array/tabela subjacente).

---

## 9. Arrays, slices e maps

**Array** tem tamanho fixo (parte do tipo):

```go
var a [3]int = [3]int{1, 2, 3}
b := [...]int{4, 5, 6} // tamanho inferido
```

**Slice** é a estrutura mais usada — visão dinâmica sobre um array:

```go
s := []int{1, 2, 3}
s = append(s, 4, 5)         // pode realocar
sub := s[1:3]               // [2 3]
length := len(s)
capacity := cap(s)

s2 := make([]int, 5)         // len=5, cap=5
s3 := make([]int, 0, 10)     // len=0, cap=10

copy(dst, src)               // copia elementos
```

Atenção: ao fatiar um slice, o novo slice **compartilha** o array subjacente. Para isolar, use `append([]T{}, src...)` ou `slices.Clone`.

**Map** é uma tabela hash:

```go
m := map[string]int{"a": 1, "b": 2}
m2 := make(map[string]int)

m["c"] = 3
v, ok := m["x"]    // ok=false se não existir
delete(m, "a")

for k, v := range m { ... }
```

Maps **não** garantem ordem de iteração e não são seguros para acesso concorrente.

---

## 10. Structs

```go
type Pessoa struct {
    Nome  string
    Idade int
}

p := Pessoa{Nome: "Ana", Idade: 30}
p2 := Pessoa{"Bia", 25}            // posicional (evite)
p3 := &Pessoa{Nome: "Carlos"}      // ponteiro

p.Idade++

// struct anônimo
ponto := struct{ X, Y int }{1, 2}

// embedding (composição)
type Funcionario struct {
    Pessoa        // tipo embutido
    Salario float64
}
f := Funcionario{Pessoa{"Ana", 30}, 5000}
fmt.Println(f.Nome) // promove campo
```

**Tags** são metadados acessíveis via reflection (usadas por `encoding/json`, etc.):

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name,omitempty"`
}
```

---

## 11. Métodos e interfaces

Método é uma função com **receiver**:

```go
type Retangulo struct{ L, A float64 }

func (r Retangulo) Area() float64 { return r.L * r.A }       // valor
func (r *Retangulo) Escalar(k float64) { r.L *= k; r.A *= k } // ponteiro
```

Use receiver de **ponteiro** quando precisa mutar o valor ou quando o struct é grande.

**Interface** define um conjunto de métodos. Implementação é **implícita** (duck typing):

```go
type Forma interface {
    Area() float64
}

func imprimir(f Forma) { fmt.Println(f.Area()) }
imprimir(Retangulo{3, 4}) // funciona porque Retangulo tem Area()
```

`interface{}` (ou `any` em Go 1.18+) aceita qualquer tipo. Recupere o tipo concreto com **type assertion** ou **type switch**:

```go
var i any = "hello"
s := i.(string)        // panic se errado
s, ok := i.(string)    // forma segura
```

---

## 12. Ponteiros

Diferentemente de C, Go **não tem aritmética de ponteiros**.

```go
x := 10
p := &x      // p é *int
*p = 20      // muta x
fmt.Println(x) // 20

n := new(int) // ponteiro para int zero
*n = 5
```

`nil` é o valor zero de ponteiros. Sempre verifique antes de desreferenciar.

---

## 13. Tratamento de erros

Erros são **valores**, não exceções. Convencionalmente o último retorno é `error`:

```go
func dividir(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("divisão por zero")
    }
    return a / b, nil
}

if v, err := dividir(10, 0); err != nil {
    log.Fatal(err)
} else {
    fmt.Println(v)
}
```

Criando erros:

```go
errors.New("mensagem")
fmt.Errorf("falha ao abrir %s: %w", path, err) // %w embrulha (wrap)

// inspecionar wraps:
errors.Is(err, os.ErrNotExist)
var pathErr *os.PathError
errors.As(err, &pathErr)
```

Tipos de erro customizados implementam a interface `error`:

```go
type ValidationError struct{ Campo string }
func (e *ValidationError) Error() string {
    return fmt.Sprintf("campo %q inválido", e.Campo)
}
```

---

## 14. Defer, panic e recover

`defer` adia uma chamada até o retorno da função. Útil para liberar recursos.

```go
func ler(path string) ([]byte, error) {
    f, err := os.Open(path)
    if err != nil { return nil, err }
    defer f.Close() // executa quando ler() retornar
    return io.ReadAll(f)
}
```

Múltiplos `defer` executam em **LIFO**. Os argumentos são avaliados imediatamente; só a chamada é adiada.

`panic` interrompe a execução; `recover` (apenas dentro de `defer`) intercepta:

```go
func safe() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recuperado de:", r)
        }
    }()
    panic("algo deu errado")
}
```

Use `panic` apenas para erros realmente irrecuperáveis ou bugs. Para fluxo normal, use `error`.

---

## 15. Goroutines e canais (concorrência)

**Goroutine** é uma função em execução concorrente leve:

```go
go func() { fmt.Println("oi") }()
```

**Channel** comunica goroutines. Mantra: *"não comunique compartilhando memória; compartilhe memória comunicando"*.

```go
ch := make(chan int)        // sem buffer
ch := make(chan int, 5)     // com buffer

ch <- 42                    // envia
v := <-ch                   // recebe
v, ok := <-ch               // ok=false se canal fechado e vazio
close(ch)                   // fecha (apenas o produtor)

for v := range ch { ... }   // lê até fechar
```

**`select`** multiplexa operações de canal:

```go
select {
case v := <-ch1:
    fmt.Println(v)
case ch2 <- 1:
    fmt.Println("enviado")
case <-time.After(time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("nada pronto")
}
```

Sincronização (`sync`):

```go
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        fmt.Println(n)
    }(i)
}
wg.Wait()

var mu sync.Mutex
mu.Lock()
// seção crítica
mu.Unlock()

var once sync.Once
once.Do(func() { fmt.Println("só uma vez") })
```

`context.Context` propaga cancelamento e deadline:

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
select {
case <-ctx.Done():
    return ctx.Err()
case res := <-ch:
    return res
}
```

---

## 16. Generics

Disponíveis a partir de Go 1.18. Parâmetros de tipo entre `[]`:

```go
func Map[T, U any](xs []T, f func(T) U) []U {
    out := make([]U, len(xs))
    for i, x := range xs { out[i] = f(x) }
    return out
}

dobros := Map([]int{1,2,3}, func(n int) int { return n*2 })
```

Restrições (constraints) limitam tipos aceitos. Pacote `constraints` e operador `~` (tipos cuja base é T):

```go
type Number interface {
    ~int | ~float64
}

func Soma[T Number](xs []T) T {
    var s T
    for _, x := range xs { s += x }
    return s
}
```

---

## 17. Testes e benchmarks

Arquivos de teste terminam em `_test.go` no mesmo pacote. Funções iniciam com `Test`, recebem `*testing.T`.

```go
package soma

import "testing"

func TestSoma(t *testing.T) {
    got := Soma(2, 3)
    want := 5
    if got != want {
        t.Errorf("Soma(2,3) = %d; quer %d", got, want)
    }
}
```

**Table-driven tests** (padrão na comunidade):

```go
func TestSoma(t *testing.T) {
    casos := []struct {
        nome       string
        a, b, want int
    }{
        {"positivos", 2, 3, 5},
        {"zeros", 0, 0, 0},
        {"negativo", -1, 1, 0},
    }
    for _, c := range casos {
        t.Run(c.nome, func(t *testing.T) {
            if got := Soma(c.a, c.b); got != c.want {
                t.Errorf("got %d, want %d", got, c.want)
            }
        })
    }
}
```

Outros utilitários:

- `t.Fatal/Fatalf` — interrompe o teste imediatamente.
- `t.Skip` — pula o teste.
- `t.Helper()` — marca função auxiliar (mostra linha do chamador).
- `t.Parallel()` — roda subtests em paralelo.
- `t.Cleanup(fn)` — registra função executada ao final.

Benchmarks (`Benchmark<Nome>`):

```go
func BenchmarkSoma(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Soma(2, 3)
    }
}
```

Roda com `go test -bench=. -benchmem`.

Exemplos executáveis aparecem na doc:

```go
func ExampleSoma() {
    fmt.Println(Soma(1, 2))
    // Output: 3
}
```

Cobertura: `go test -cover` ou `go test -coverprofile=cover.out`.

---

## 18. Pacotes da biblioteca padrão mais usados

| Pacote | Para que serve |
|---|---|
| `fmt` | Entrada/saída formatada (`Println`, `Printf`, `Errorf`, `Sprintf`). |
| `strings` | Manipulação de strings. |
| `strconv` | Conversão entre strings e tipos numéricos/bool. |
| `bytes` | Buffers e manipulação de `[]byte`. |
| `errors` | Criar e inspecionar erros (`Is`, `As`, `Unwrap`, `Join`). |
| `os` | Sistema operacional: arquivos, args, env (`os.Open`, `os.Args`, `os.Getenv`). |
| `io` / `io/fs` | Abstrações de leitura/escrita (`Reader`, `Writer`, `Copy`, `ReadAll`). |
| `bufio` | I/O bufferizado (`Scanner`, `NewReader`). |
| `path/filepath` | Manipulação de caminhos do SO. |
| `time` | Datas, durações, timers (`time.Now`, `time.Sleep`, `time.Since`). |
| `math`, `math/rand` | Funções matemáticas e RNG. |
| `sort` / `slices` | Ordenação e utilidades de slices. |
| `encoding/json` | Marshal/Unmarshal JSON. |
| `net/http` | Cliente e servidor HTTP. |
| `database/sql` | Acesso a bancos relacionais (com drivers). |
| `context` | Cancelamento, deadlines, valores por requisição. |
| `sync`, `sync/atomic` | Primitivas de concorrência. |
| `log`, `log/slog` | Logging (texto e estruturado). |
| `regexp` | Expressões regulares. |
| `flag` | Parser simples de flags de CLI. |
| `testing` | Framework de testes/benchmarks/exemplos. |

Exemplo: servidor HTTP em ~10 linhas

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "olá", r.URL.Path)
    })
    http.ListenAndServe(":8080", nil)
}
```

Exemplo: JSON

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

data, _ := json.Marshal(User{1, "Ana"})        // codifica
var u User
_ = json.Unmarshal(data, &u)                   // decodifica
```

---

## 19. Convenções e boas práticas

- **Formate** com `gofmt`/`go fmt` sempre (sem discussão de estilo).
- Nomes curtos e idiomáticos: `i`, `r`, `ctx`, `err`, `buf`. Quanto **menor o escopo, mais curto o nome**.
- Nomes de pacotes: minúsculos, curtos, sem underscores; **não repita** o nome do pacote nos identificadores (`http.Server` em vez de `http.HTTPServer`).
- Erros começam com letra **minúscula** e não terminam em ponto: `errors.New("conexão recusada")`.
- Trate todo `error` retornado. Não ignore com `_` sem motivo.
- Prefira **early return** a aninhamento profundo de `if/else`.
- Funções/arquivos curtos; um arquivo por tópico.
- Use ponteiros somente quando precisar (mutação ou structs grandes).
- Evite alocações desnecessárias; pré-aloque slices com `make([]T, 0, n)` quando souber o tamanho.
- Documente identificadores exportados — comentário começando com o **próprio nome**:

```go
// Soma retorna a + b.
func Soma(a, b int) int { return a + b }
```

- Linters recomendados: `go vet`, `staticcheck`, `golangci-lint`.
- Prefira **composição** (embedding) a "herança".
- Interfaces devem ser **pequenas** e definidas no **lado de quem consome**, não de quem implementa.
- Concorrência: comunique por canais; use mutex quando o estado é local; sempre tenha um plano para encerrar goroutines (via `context` ou canal de "done") para não vazar.

---

## Referências oficiais

- [Tour of Go](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Documentação da std](https://pkg.go.dev/std)
- [Go Proverbs](https://go-proverbs.github.io/)
