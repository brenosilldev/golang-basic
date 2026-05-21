package fundacao

import (
	"fmt"
)

func ForTeste() {

	for i := 0; i < 10; i++ {
		print("\n", i)
	}

	num := []string{"um", "dois", "3"}

	for k, v := range num {
		fmt.Println(k, v)
	}

	i := 0

	for i < 10 {
		fmt.Println(i)
		i++
	}

}
