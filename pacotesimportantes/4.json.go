package pacotesimportantes

import (
	"encoding/json"
	"fmt"
	"os"
)

// Definindo o json
type Conta struct {
	Saldo  int `json:"saldo"`
	Numero int `json:"nmr"`
}

func JsonTes() {
	conta := Conta{Numero: 1, Saldo: 100}

	res, err := json.Marshal(conta)

	if err != nil {
		println(err)
	}

	fmt.Println("Res: ", res)
	fmt.Println("Res: ", string(res))

	// Retorna direto no terminal = os.Stdout
	json.NewEncoder(os.Stdout).Encode(conta)

	jsonPuro := []byte(`{"nmr":2,"saldo":200}`)
	var contaX Conta
	err = json.Unmarshal(jsonPuro, &contaX)

	if err != nil {
		println(err)
	}

	println("Saldo: ", contaX.Numero)

}
