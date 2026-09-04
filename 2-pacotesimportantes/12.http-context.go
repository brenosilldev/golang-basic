package pacotesimportantes

import (
	"context"
	"fmt"
	"time"
)

//Explicação do pacote context: O pacote context fornece uma maneira de passar informações entre goroutines, como deadlines, cancelamentos e valores de contexto. Ele é amplamente utilizado em aplicações web para gerenciar requisições HTTP e controlar o tempo de execução de operações assíncronas.
//Context = Context é um pacote que permite passar informações entre goroutines, como deadlines, cancelamentos e valores de contexto. Ele é amplamente utilizado em aplicações web para gerenciar requisições HTTP e controlar o tempo de execução de operações assíncronas.

func HttpContextTest() {

	ctx := context.Background()

	// Criando um contexto com timeout de 2 segundos
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

	for i := 1; i < 5; i++ {
		go func(i int) {
			select {
			case <-time.After(2 * time.Second):
				fmt.Println("Goroutine", i, "concluída com sucesso!")
			case <-ctx.Done():
				fmt.Println("Goroutine", i, "cancelada devido ao timeout.")
			}
		}(i)
	}

	// Aguardando o término das goroutines
	time.Sleep(3 * time.Second)
	fmt.Println("Fim do programa.")

	defer cancel()

}
