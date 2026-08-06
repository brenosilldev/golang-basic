package fundacao

import "fmt"

/*
	Slices são estruturas de dados que permitem armazenar uma sequência de elementos do mesmo tipo

	make = cria uma slice com tamanho e capacidade definidos


*/

func SliceTeste() {

	slice := []int{2, 4, 6, 8, 10}
	fmt.Println("----- Slice (literal) -----")
	fmt.Printf("len=%d cap=%d %v\n", len(slice), cap(slice), slice)

	fmt.Println("----- Slice com make -----")
	slice2 := make([]int, 3, 5)
	fmt.Printf("len=%d cap=%d %v\n", len(slice2), cap(slice2), slice2)

	fmt.Println("----- Rebanadas (slicing) -----")
	fmt.Printf("slice[:] = %v\n", slice[:])     // visão do slice inteiro
	fmt.Printf("slice[1:] = %v\n", slice[1:])   // do índice 1 até o fim
	fmt.Printf("slice[2:] = %v\n", slice[2:])   // do índice 2 até o fim
	fmt.Printf("slice[1:3] = %v\n", slice[1:3]) // índices 1 e 2 (intervalo [1,3))
	fmt.Printf("slice[1:4] = %v\n", slice[1:4]) // índices 1, 2 e 3 (intervalo [1,4))

	fmt.Println("----- Slice com append -----")
	fmt.Printf("slice = %v  cap=%d\n", slice, cap(slice))
	slice3 := append(slice, 12)
	fmt.Printf("slice3 = %v  cap=%d\n", slice3, cap(slice3))

}
