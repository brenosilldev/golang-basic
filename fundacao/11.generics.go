package fundacao

import "fmt"

/*
 */
func Soma[T int | float64](m map[string]T) T {
	var soma T
	for _, v := range m {
		soma += v
	}

	return soma
}

func GenericsTeste() {

	fmt.Println(Soma(map[string]int{"Breno": 2000, "Hose": 2000}))

}
