package pacotesimportantes

import (
	"encoding/json"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ViaCe2 struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
}

func HttpTest2() {

	http.HandleFunc("/cep", BuscarCephandler)

	http.ListenAndServe(":8080", nil)

}

func BuscarCephandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/cep" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	cep := r.URL.Query().Get("cep")

	if cep == "" {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Bad Request",
			"message": "O parâmetro 'cep' é obrigatório.",
		})

		return
	}

	viaCep, err := BuscarCep(cep)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "CEP não encontrado",
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusProcessing)
	json.NewEncoder(w).Encode(viaCep)

}
 
func BuscarCep(cep string) (*ViaCe2, error) {

	resp, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err

	}

	var viaCep ViaCe2

	err = json.Unmarshal(body, &viaCep)

	if err != nil {
		return nil, err
	}

	return &viaCep, nil
}
