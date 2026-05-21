package pacotesimportantes

import (
	"fmt"
	"os"
)

func PacotesTest() {

	f, err := os.Create("arquivo.txt")

	if err != nil {
		panic(err)
	}

	tamanho, err := f.Write([]byte("Ola mundo"))

	fmt.Printf("ARquivo criado %d \n", tamanho)

	f.Close()

	arq, err := os.ReadFile("arquivo.txt")

	if err != nil {
		panic(err)
	}

	fmt.Println(string(arq))

}
