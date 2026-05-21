package fundacao

import "fmt"

// constrait para generics
type Number interface {
	~int | float64
}
type Mynumber int

func Soma[T Number](m map[string]T) T {

	var soma T
	for i, v := range m {
		fmt.Printf("i: %s \n", i)
		soma += v
	}

	return soma
}

// funcção que compara tipos
func Compara[T comparable](a, b T) bool {
	if a == b {
		return true
	}

	return false
}

func GenericsTeste() {

	fmt.Println(Soma(map[string]Mynumber{"Breno": 2000, "Hose": 2000}))
	fmt.Println(Compara("a", "b"))

}
