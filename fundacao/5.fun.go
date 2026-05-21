package fundacao

import (
	"errors"
	"fmt"
)

/*
	Função é um bloco de código que pode ser chamado

	func = define uma função
	func nome(parametros) retorno {
		bloco de código
	}
*/

func FuncaoTeste() {

	fmt.Println("Função Teste")
	valor, err := sum(20, 19)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(valor)
	}

	valor2 := sum2(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	fmt.Printf("soma dos pares: %d\n", valor2)

}

func sum(a, b int) (int, error) {
	if a+b >= 20 {
		return 0, errors.New("A soma maior foi maior que 20")
	} else {
		return a + b, nil
	}
}

//Função variadicas

func sum2(numeros ...int) int {
	total := 0
	for _, numeros := range numeros {
		if numeros%2 == 0 {
			total += numeros
		}
	}

	return total
}
