package fundacao

import (
	"fmt"
	"net/mail"
)

type User struct {
	Nome  string
	Email string
	Ativo bool
	Endereco
}

type Endereco struct {
	Rua    string
	Numero int
}

func (u User) Desativar() {
	u.Ativo = false
	fmt.Printf("o cliente %s foi desativado \n", u.Nome)
}

func (u User) validarEmail() bool {
	_, err := mail.ParseAddress(u.Email)

	return err == nil
}

func StructsTeste() {

	breno := User{
		Nome:  "Breno Silva de Oliveira",
		Email: "brenosill@hotmail.com",
		Endereco: Endereco{
			Rua:    "Travess Luciano torrente",
			Numero: 24,
		},
	}

	breno.Desativar()
	fmt.Printf("email valido: %t\n", breno.validarEmail())
}
