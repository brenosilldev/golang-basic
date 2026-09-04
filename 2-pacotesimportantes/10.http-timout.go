package pacotesimportantes

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"bytes"
)

func HttpTimeoutTest() {

	client := http.Client{}

	jsv := bytes.NewBuffer([]byte(`{"name": "Breno Silva"}`))
	// Fazendo uma requisição GET para um site que demora para responder
	// resp, err := client.Get("https://google.com.br/search")

	// Fazendo uma requisição POST para um site que demora para responder
	resp, err := client.Post("https://google.com.br", "application/json", jsv)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return
	}

	defer resp.Body.Close()

	// Lendo o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return
	}

	fmt.Println("Resposta recebida com sucesso!")
	fmt.Println(string(body))
}
