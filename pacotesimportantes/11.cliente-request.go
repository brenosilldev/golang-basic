package pacotesimportantes

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func ClientRequestTest() {
	c := http.Client{}

	// Fazendo uma requisição GET para o site do Google
	req, err := http.NewRequest("GET", "https://www.google.com.br", nil)
	if err != nil {
		panic(err)
	}

	// Adicionando um cabeçalho à requisição
	req.Header.Set("Accept", "application/json")

	// Fazendo a requisição e obtendo a resposta
	resp, err := c.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Lendo o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Resposta recebida com sucesso!")
	fmt.Println(string(body))

	fmt.Println("-------------------")
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Headers:", resp.Header)
}
