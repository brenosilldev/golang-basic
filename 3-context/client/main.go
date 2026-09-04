package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Cria um contexto com timeout de 500ms

	defer cancel() // Garante que o cancelamento do contexto seja chamado ao final da função
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/teste", nil)

	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req) // Executa a requisição HTTP com o contexto

	if err != nil {
		log.Println("Erro na requisição:", err)
		return
	}

	defer resp.Body.Close() // Garante que o corpo da resposta seja fechado ao final da função

	io.Copy(os.Stdout, resp.Body) // Lê o corpo da resposta e descarta o conteúdo
	log.Println("Resposta recebida com status:", resp.Status)

}
