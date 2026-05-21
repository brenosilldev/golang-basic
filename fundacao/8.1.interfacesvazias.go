package fundacao

import "fmt"

func imprimir(valor interface{}) {
	fmt.Println(valor)
}

func InterfaceVaziasTestes() {

	var x interface{} = "Breno"

	imprimir(10)
	imprimir("Go")
	imprimir(true)
	imprimir(x)
}
