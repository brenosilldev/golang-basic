package fundacao

import "fmt"

/*
	PONTEIROS em Go

	Toda variável ocupa um espaço na memória com um endereço único.
	Um ponteiro é simplesmente uma variável que guarda esse endereço.

	Por que usar ponteiros?
	  - Sem ponteiro: a função recebe uma CÓPIA do valor. Alterar a cópia não muda o original.
	  - Com ponteiro: a função recebe o ENDEREÇO. Alterar via ponteiro muda o original.

	Os dois símbolos:

	  &variavel   →  "me dá o endereço de variavel"  (cria o ponteiro)
	  *ponteiro   →  "me dá o valor que está nesse endereço"  (desreferencia)

	  *int        →  tipo que significa "ponteiro para um int"  (no tipo, não na expressão)

	Mapa mental:
	  n    = 10          (o valor em si)
	  &n   = 0xc000...   (endereço onde n mora)
	  p    = &n          (p guarda o endereço)
	  *p   = 10          (lê o valor no endereço que p aponta — mesmo que n)
	  *p = 99            (escreve 99 no endereço — altera n diretamente)
*/

// Pessoa é usada nos exemplos de ponteiro para struct.
type Pessoa struct {
	Nome  string
	Email string
}

func PonteirosTestes() {
	fmt.Println("=== Ponteiros ===")
	fmt.Println()

	// --- 1. CÓPIA vs PONTEIRO ---
	// Quando passamos n diretamente, a função recebe uma cópia.
	// A variável original não é tocada.
	n := 10
	fmt.Printf("antes: n = %d\n", n)

	alterarCopia(n)
	fmt.Printf("depois de alterarCopia:    n = %d  (não mudou — era só uma cópia)\n", n)

	// &n envia o endereço de n. A função pode alterar o original.
	alterarOriginal(&n)
	fmt.Printf("depois de alterarOriginal: n = %d  (mudou — recebeu o endereço)\n", n)

	fmt.Println()

	// --- 2. LENDO E ESCREVENDO PELO PONTEIRO ---
	a := 42
	p := &a // p é do tipo *int; guarda o endereço de a

	fmt.Printf("a    = %d\n", a)
	fmt.Printf("&a   = %p  (endereço de a)\n", &a)
	fmt.Printf("p    = %p  (p guarda esse mesmo endereço)\n", p)
	fmt.Printf("*p   = %d  (valor no endereço que p aponta — igual a a)\n", *p)

	*p = 100 // escreve 100 diretamente no endereço de a
	fmt.Printf("após *p = 100 →  a = %d  (a foi alterado via ponteiro)\n", a)

	fmt.Println()

	// --- 3. PONTEIRO PARA STRUCT ---
	// Acessar campos de um ponteiro para struct com "." funciona normalmente.
	// Go faz a desreferência automática: pessoa.Nome == (*pessoa).Nome
	pessoa := &Pessoa{Nome: "Maria", Email: "maria@email.com"}
	fmt.Printf("pessoa.Nome  = %s\n", pessoa.Nome)
	fmt.Printf("pessoa.Email = %s\n", pessoa.Email)

	atualizarEmail(pessoa, "novo@email.com")
	fmt.Printf("após atualizarEmail: %s\n", pessoa.Email)

	fmt.Println()

	// --- 4. DOIS PONTEIROS, MESMA MEMÓRIA ---
	// x e y apontam para a mesma variável. Alterar via y também afeta x.
	valor := 5
	px := &valor
	py := &valor // mesmo endereço

	*px = 50
	fmt.Printf("*py = %d *px %d  (y enxerga a mudança feita por x: mesma memória)\n", *py,*py	)
}

// alterarCopia recebe uma cópia de x. Alterar x aqui não tem efeito fora.
func alterarCopia(x int) {
	x = 99
}

// alterarOriginal recebe o endereço (*int). *x = 99 altera a variável original.
func alterarOriginal(x *int) {
	*x = 99
}

// atualizarEmail recebe um ponteiro para Pessoa e altera o campo diretamente.
// Se recebesse Pessoa por valor, a alteração seria descartada ao retornar.
func atualizarEmail(p *Pessoa, novoEmail string) {
	p.Email = novoEmail // Go desreferencia automaticamente: p.Email == (*p).Email
}
