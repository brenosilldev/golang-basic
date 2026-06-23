package pacotesimportantes

import "fmt"

/*
Defer atrasa chamada desejada
*/
func DeferTest() {

	fmt.Println("Ola mundo")

	defer fmt.Println("Ultima chamada ")

	fmt.Println("No meio do projeto")

}
