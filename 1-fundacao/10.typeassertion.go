package fundacao

import "fmt"

func TypeAssertionTeste() {
	var minhavar interface{} = "Wesley"
	res, ok := minhavar.(int)

	fmt.Printf("o Valor de res %v , e o resultado ok é %v \n", res, ok)
}
