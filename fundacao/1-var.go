package fundacao

import "fmt"

// constantes = constantes não podem ser alteradas
const (
	Versao = 1.0
)

// variáveis = variáveis podem ser alteradas
var (
	Fundacao = "Fundacao"
)

// tipos = tipos são definições de tipos
type ID int

func VarTeste(nome1 string) {

	var id ID = 1
	fmt.Println("ID", id)
	nome := nome1
	fmt.Println("Nome", nome)
	fmt.Println("Fundacao Atualizada", Versao)
	fmt.Println("Fundacao", Fundacao)

	fmt.Println("--------------------------------")
	fmt.Printf("O tipo de id é %T \n", id)
	fmt.Printf("O tipo de nome é %T \n", nome)
	fmt.Printf("O tipo de Versao é %T \n", Versao)
	fmt.Printf("O tipo de Fundacao é %T \n", Fundacao)
}
