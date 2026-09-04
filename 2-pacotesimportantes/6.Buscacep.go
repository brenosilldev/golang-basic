package pacotesimportantes

// Importação dos pacotes necessários.
import (
	"encoding/json" // Converte JSON para structs Go e vice-versa.
	"fmt"           // Funções de formatação e impressão.
	"io"            // Utilizado para ler o corpo da resposta HTTP.
	"net/http"      // Permite fazer requisições HTTP.
	"os"            // Manipulação de arquivos e argumentos do terminal.
)

/*
	Forma de executar:

	go run main.go 02873420

	O CEP será capturado através de os.Args.
*/

// Struct que representa exatamente o JSON retornado pela API ViaCEP.
// As tags (`json:"..."`) indicam qual campo do JSON corresponde
// a cada atributo da struct.
type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

// Função responsável por buscar um ou mais CEPs informados
// pela linha de comando.
func BuscaCepTest() {

	// os.Args contém todos os argumentos do programa.
	//
	// Exemplo:
	// go run main.go 02873420 01001000
	//
	// os.Args[0] -> nome do programa
	// os.Args[1] -> 02873420
	// os.Args[2] -> 01001000
	//
	// O laço percorre apenas os CEPs informados.
	for _, cep := range os.Args[1:] {

		// Faz uma requisição GET para a API do ViaCEP.
		req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")

		// Verifica se ocorreu algum erro durante a requisição.
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
			continue
		}

		// Garante que o corpo da resposta será fechado
		// quando a função terminar.
		defer req.Body.Close()

		// Lê todo o conteúdo retornado pela API.
		// O resultado será um slice de bytes ([]byte).
		res, err := io.ReadAll(req.Body)

		// Verifica se ocorreu erro na leitura.
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
			continue
		}

		// Cria uma variável do tipo ViaCep.
		// Ela será preenchida com os dados do JSON.
		var data ViaCep

		// Converte o JSON recebido para a struct.
		err = json.Unmarshal(res, &data)

		// Verifica se houve erro na conversão.
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
			continue
		}

		// Neste momento, todos os campos da struct já estão preenchidos.
		//
		// Exemplo:
		// data.Cep
		// data.Logradouro
		// data.Bairro
		// data.Localidade
		// data.Uf

		// Cria (ou sobrescreve) um arquivo chamado "cep.json".
		file, err := os.Create("cep.json")

		// Verifica se houve erro ao criar o arquivo.
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
			continue
		}

		// Escreve exatamente o JSON retornado pela API
		// dentro do arquivo.
		_, err = file.WriteString(string(res))

		// Verifica se houve erro na escrita.
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		}

		// Fecha o arquivo ao final da função.
		defer file.Close()

		// Se desejar utilizar os dados da struct:
		//
		// fmt.Println(data.Cep)
		// fmt.Println(data.Logradouro)
		// fmt.Println(data.Bairro)
		// fmt.Println(data.Localidade)
	}
}
