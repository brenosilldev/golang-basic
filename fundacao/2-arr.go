package fundacao

import "fmt"

func ArrTeste() {

	arr := [3]int{1, 2, 3} // array de 3 elementos
	fmt.Println("Arr", arr)
	fmt.Println("--------------------------------")
	fmt.Printf("O tipo de arr é %T \n", arr)
	fmt.Println("O tamanho do array é", len(arr))
	fmt.Println("O primeiro elemento do array é", arr[0])
	fmt.Println("O último elemento do array é", arr[len(arr)-1])
	fmt.Println("O array completo é", arr)

	for i, v := range arr {
		fmt.Println("O indice do array é", i, "e o valor é", v)
	}
}
