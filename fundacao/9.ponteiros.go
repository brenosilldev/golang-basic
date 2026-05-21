package fundacao

/*
	Ponteiros — os símbolos & e *

	&  = "endereço de"
	     Pega o ponteiro para uma variável.
	     Exemplo: p := &n  →  p guarda ONDE n está na memória.

	*  tem dois usos:

	     1) No TIPO  →  significa "ponteiro para"
	        Exemplo: *int  →  variável que aponta para um int
	        Exemplo: *User →  variável que aponta para um User

	     2) Na EXPRESSÃO  →  significa "valor apontado" (desreferenciar)
	        Exemplo: *p  →  lê ou altera o valor que p aponta
	        Exemplo: *p = 99  →  altera o int original, não a cópia

	Resumo rápido:
	     n   →  o valor (10)
	     &n  →  o endereço de n (ponteiro *int)
	     p   →  o endereço guardado
	     *p  →  o valor que está nesse endereço
*/

type Pessoa struct {
	Nome  string
	Email string
}

func PonteirosTestes() {

	// explicarOperadores()

	// fmt.Println("--- ponteiros: int ---")

	// n := 10
	// fmt.Printf("antes: n = %d\n", n)

	// alterarCopia(n)
	// fmt.Printf("depois de alterarCopia (cópia): n = %d\n", n)

	// alterarOriginal(&n) // &n = passa o endereço de n para a função alterar o original
	// fmt.Printf("depois de alterarOriginal (ponteiro): n = %d\n", n)

	// a := 10
	// fmt.Println("valor de A :", a)
	// ponteiro := &a
	// fmt.Printf("Valor ponteiro na memoria: %p \n", ponteiro)

	// b := *ponteiro
	// fmt.Println("Valor do ponteiro: ", b)
	// c := b + 10 // Recebe o valor de ponteiro + 10
	// fmt.Println("valor de ponteiro somado: ", c)
	// fmt.Println("valor de ponteiro somado: ", c+10)

	// var1 := 10
	// var2 := 20
	// println(soma(&var1, &var2))
	// println(var1)
	// println(var2)

}

// func soma(a, b *int) int {

// 	fmt.Println(*a)
// 	fmt.Println(*b)

// 	*a = 20
// 	return *a + *b
// }

// func explicarOperadores() {
// 	fmt.Println("--- & e * na prática ---")

// 	n := 10

// 	// &n → ponteiro para n (tipo *int)
// 	p := &n
// 	fmt.Printf("n = %d (valor)\n", n)
// 	fmt.Printf("&n → p = %p (endereço / ponteiro)\n", p)
// 	fmt.Printf("*p = %d (valor no endereço que p guarda)\n", *p)

// 	// *p altera o original, porque acessa o mesmo lugar que n
// 	*p = 25
// 	fmt.Printf("depois de *p = 25 → n = %d\n", n)

// 	// Em parâmetros: *int no tipo = "recebo um ponteiro"
// 	// Na chamada: &n = "envio o endereço de n"
// 	fmt.Println()
// 	fmt.Println("&  →  pega o endereço   (n → &n)")
// 	fmt.Println("*  →  no tipo: ponteiro | na expressão: valor apontado")
// 	fmt.Println()
// }

// func alterarCopia(x int) {
// 	x = 99 // só altera a cópia local
// }

// func alterarOriginal(x *int) {
// 	*x = 99 // *x = altera o int que x aponta (o original)
// }
