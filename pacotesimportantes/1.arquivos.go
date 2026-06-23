package pacotesimportantes

import (
	"bufio"
	"fmt"
	"os"
)

func PacotesTest() {

	f, err := os.Create("arquivo.txt")

	if err != nil {
		panic(err)
	}

	tamanho, err := f.Write([]byte("Ola mundo"))

	fmt.Printf("Arquivo criado %d \n", tamanho)

	f.Close()

	arq, err := os.ReadFile("arquivo.txt")

	if err != nil {
		panic(err)
	}

	fmt.Println(string(arq))

	// Leitura de pouco em pouco abrindo o arquivo
	arq2, err := os.Open("arquivo.txt")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(arq2)
	buffer := make([]byte, 3)

	for {
		n, err := reader.Read(buffer)

		if err != nil {
			break

		}

		fmt.Println(string(buffer[:n]))
	}


	//Remover 
	err = os.Remove("arquivo.txt")

	if err != nil {
		panic(err)
	}

}
