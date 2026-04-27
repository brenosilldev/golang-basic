package modulo3

import (
	"fmt"
	"github.com/badoux/checkmail"
)

func ControleFluxo() {
	fmt.Println("=== CONTROLE DE FLUXO ===")

	err := checkmail.ValidateFormat("teste@teste.com")
	
	if err != nil {
		fmt.Println("Email inválido")
	} else {
		fmt.Println("Email válido")
	}


}