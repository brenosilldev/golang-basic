package fundacao

import "fmt"

func ClosuresTeste() {

	// função clouseres
	total := func() int {
		return sum2(1, 2, 3, 5, 6, 7, 8, 9, 0, 11, 29, 10)
	}()

	fmt.Printf("valor soma: %d \n", total)

}
