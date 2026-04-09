package modulo2

import (
	"fmt"
	"math"
	// "strings"
)

func VarExample() {

	// Declaração de variáveis com tipo explícito
	// var nome string = "Breno"
	// var idade int = 26
	// var altura float64 = 1.75
	// var isDeveloper bool = true
	// var linguagens []string = []string{"Go", "Python", "JavaScript"}
	// var skills map[string]string = map[string]string{"Go": "Intermediário", "Python": "Básico", "JavaScript": "Básico"}
	// var endereco string = "Rua dos Bobos, 0"
	// var cidade string = "São Paulo"
	// var estado string = "SP"
	// var pais string = "Brasil"

	// // Declaração de variáveis com tipo implícito
	// nome := "Breno"
	// idade := 26
	// altura := 1.75
	// isDeveloper := true
	// linguagens := []string{"Go", "Python", "JavaScript"}
	// skills := map[string]string{"Go": "Intermediário", "Python": "Básico", "JavaScript": "Básico"}
	// endereco := "Rua dos Bobos, 0"
	// cidade := "São Paulo"
	// estado := "SP"
	// pais := "Brasil"


	// // Declaração de variáveis com tipo explícito e múltiplas variáveis
	// var (
	// 	nome string = "Breno"
	// 	idade int = 26
	// 	altura float64 = 1.75
	// 	isDeveloper bool = true
	// 	linguagens []string = []string{"Go", "Python", "JavaScript"}
	// 	skills map[string]string = map[string]string{"Go": "Intermediário", "Python": "Básico", "JavaScript": "Básico"}
	// 	endereco string = "Rua dos Bobos, 0"
	// 	cidade string = "São Paulo"
	// 	estado string = "SP"
	// 	pais string = "Brasil"
	//  	byte1 byte = 123
	// )

	// // fmt.Println(nome, idade, altura, isDeveloper, linguagens, skills, endereco, cidade, estado, pais, byte1)

	// var sb strings.Builder
	// for i := 0; i < 1000; i++ {
	// 	sb.WriteString("a"  + " \n") // eficiente — aloca uma vez
	// }
	// resultado := sb.String()

	// fmt.Println(resultado)


	// type Permissao int
	// const (
	// 	Leitura   Permissao = 1 << iota // 1 (0001)
	// 	Escrita                         // 2 (0010)
	// 	Execucao                        // 4 (0100)
	// )
	// fmt.Println(Leitura, Escrita, Execucao)


	fmt.Println("=== VARIÁVEIS ===")

	// --- Declaração com var ---
	// var nome string = "Go"
	// var idade int = 15 // Go foi lançada em 2009
	// var ativa bool = true

	// fmt.Printf("Linguagem: %s, Idade: %d, Ativa: %t\n", nome, idade, ativa)


	// cidade := "São Paulo" // string
	// temperatura := 28.5   // float64
	// populacao := 12000000 // int

	// fmt.Printf("Cidade: %s, Temp: %.1f°C, Pop: %d\n", cidade, temperatura, populacao)


	// // --- Múltiplas variáveis de uma vez ---
	// var (
	// 	x int     = 10
	// 	y float64 = 20.5
	// 	z string  = "hello"
	// )
	// // fmt.Printf("x=%d, y=%.1f, z=%s\n", x, y, z)

	// // ==========================================================
	// // fmt.Println("\n=== ZERO VALUES ===")

	// // Em Go, variáveis NUNCA são "undefined" ou "null"
	// // Cada tipo tem um "zero value" padrão
	// var intZero int
	// var floatZero float64
	// var stringZero string
	// var boolZero bool

	// // fmt.Printf("int: %d\n", intZero)         // 0
	// // fmt.Printf("float64: %f\n", floatZero)   // 0.000000
	// // fmt.Printf("string: '%s'\n", stringZero) // "" (vazio)
	// // fmt.Printf("bool: %t\n", boolZero)       // false


	// // int e uint variam de tamanho conforme a plataforma (32 ou 64 bits)
	// var i8 int8 = 127     // -128 a 127
	// var i16 int16 = 32767 // -32768 a 32767
	// var i32 int32 = 100000
	// var i64 int64 = 9999999999

	// fmt.Printf("int8: %d, int16: %d, int32: %d, int64: %d\n", i8, i16, i32, i64)

	
	// const Pi = 3.14159
	// const AppNome = "MeuApp"
	// fmt.Printf("Pi: %f, App: %s\n", Pi, AppNome)

	// // iota → auto-incrementa dentro de um bloco const
	// const (
	// 	Domingo = iota // 0 (auto-incrementa)
	// 	Segunda        // 1
	// 	Terca          // 2
	// 	Quarta         // 3 (auto-incrementa)
	// 	Quinta         // 4 (auto-incrementa)
	// 	Sexta          // 5 (auto-incrementa)
	// 	Sabado         // 6 (auto-incrementa)
	// )
	// 	fmt.Printf("Domingo=%d, Quinta=%d, Sexta=%d\n", Domingo, Quinta, Sexta)


	// type Permissao int
	// const (
	// 	Leitura   Permissao = 1 << iota // 1 (0001)
	// 	Escrita                       // 2 (0010)
	// 	Execucao                        // 4 (0100)
	// )

	// fmt.Printf("Leitura=%d, Escrita=%d, Execucao=%d\n", Leitura, Escrita, Execucao)
	// permissao := Leitura | Escrita // 3 (011) 
	// fmt.Printf("Leitura+Escrita = %d (%b em binário)\n", permissao, permissao)


	// // Go NÃO faz conversão implícita. Você precisa ser explícito.
	// inteiro := 42
	// decimal := float64(inteiro) // int → float64
	// fmt.Printf("int %d → float64 %.2f\n", inteiro, decimal)

	// nota := 9.7
	// notaInteira := int(nota) // float64 → int (trunca, NÃO arredonda!)
	// fmt.Printf("float64 %.1f → int %d (truncado!)\n", nota, notaInteira)

}




// ============================================================================
// EXERCÍCIO 2 — Variáveis e Tipos
// ============================================================================
//
// Exercício 2.1 — Calculadora de IMC
// Crie variáveis para peso (float64) e altura (float64).
// Calcule o IMC (peso / altura²) e imprima o resultado formatado.
// Classifique: Abaixo do peso (<18.5), Normal (18.5-24.9),
//              Sobrepeso (25-29.9), Obesidade (>=30)
//
// Exercício 2.2 — Conversor de Temperatura
// Crie uma variável celsius (float64) com um valor.
// Converta para Fahrenheit (F = C*9/5 + 32) e Kelvin (K = C + 273.15).
// Imprima os 3 valores formatados com 1 casa decimal.
//
// Exercício 2.3 — Explorar Tipos
// Declare variáveis de TODOS os tipos primitivos que aprendeu:
// int, int8, int16, int32, int64, float32, float64, string, bool, byte, rune
// Imprima cada uma com seu tipo usando %T e seu valor usando %v.
//
// Exercício 2.4 — Constantes com iota
// Crie constantes para os meses do ano usando iota (Janeiro=1, Fevereiro=2...).
// Dica: use iota + 1 para começar de 1.
// Imprima alguns meses para validar.
//
// ============================================================================

func Exercicio1() {
	// TODO: Exercício 2.1 — Calculadora de IMC
	fmt.Println("=== EXERCÍCIO 2.1 ===")

	var peso float64 = 70.5
	var altura float64 = 1.75

	fmt.Printf("Peso: %.2f, Altura: %.2f\n", peso, altura)
	fmt.Println("Calculando IMC...")
	imc := peso / math.Pow(altura, 2) // math.Pow é uma função que calcula a potência de um número
	fmt.Printf("IMC calculado: %.2f\n", imc)

	if imc < 18.5 {
		fmt.Println("Abaixo do peso")
	} else if imc < 24.9 {
		fmt.Println("Normal")
	} else if imc < 29.9 {
		fmt.Println("Sobrepeso")
	} else {
		fmt.Println("Obesidade")
	}

	// TODO: Exercício 2.2 — Conversor de Temperatura

	fmt.Println("=== EXERCÍCIO 2.2 ===")

	var celsius float64 = 35.5

	var fahrenheit float64 = celsius * 9/5 + 32
	var kelvin float64 = celsius + 273.15

	fmt.Printf("Temperatura em Celsius: %.1f°C, Fahrenheit: %.1f°F, Kelvin: %.1fK\n", celsius, fahrenheit, kelvin)

	// TODO: Exercício 2.3 — Explorar Tipos

	fmt.Println("=== EXERCÍCIO 2.3 ===")

	var (
		intZero int
		int8Zero int8
		int16Zero int16
		int32Zero int32
		int64Zero int64
		float32Zero float32
		float64Zero float64
		stringZero string
		boolZero bool
		byteZero byte
		runeZero rune
	)

	fmt.Printf("Tipo: %T, Valor: %v\n", intZero, intZero)
	fmt.Printf("Tipo: %T, Valor: %v\n", int8Zero, int8Zero)
	fmt.Printf("Tipo: %T, Valor: %v\n", int16Zero, int16Zero)
	fmt.Printf("Tipo: %T, Valor: %v\n", int32Zero, int32Zero)
	fmt.Printf("Tipo: %T, Valor: %v\n", int64Zero, int64Zero)
	fmt.Printf("Tipo: %T, Valor: %v\n", float32Zero, float32Zero)
	fmt.Printf("Tipo: %T, Valor: %v\n", float64Zero, float64Zero)
	fmt.Printf("Tipo: %T, Valor: %v\n", stringZero, stringZero)
	fmt.Printf("Tipo: %T, Valor: %v\n", boolZero, boolZero)
	fmt.Printf("Tipo: %T, Valor: %v\n", byteZero, byteZero)
	fmt.Printf("Tipo: %T, Valor: %v\n", runeZero, runeZero)

	// TODO: Exercício 2.4 — Constantes com iota

	fmt.Println("=== EXERCÍCIO 2.4 ===")

	const (
		janeiro = iota + 1
		fevereiro
		marco
		abril
		maio
		junho
		julho
		agosto
		setembro
		outubro
		novembro
		dezembro
	)

	fmt.Printf("Janeiro: %d \n" + "Fevereiro: %d \n" + "Março: %d \n" + "Abril: %d \n" + "Maio: %d \n" + "Junho: %d \n" + "Julho: %d \n" + "Agosto: %d \n" + "Setembro: %d \n" + "Outubro: %d \n" + "Novembro: %d \n" + "Dezembro: %d\n", janeiro, fevereiro, marco, abril, maio, junho, julho, agosto, setembro, outubro, novembro, dezembro)
}


