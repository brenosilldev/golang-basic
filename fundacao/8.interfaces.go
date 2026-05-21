package fundacao

import "fmt"

/*
	Interface = contrato: define quais métodos um tipo precisa ter
	para ser tratado como aquele tipo.

	Em Go, um tipo implementa uma interface implicitamente —
	basta ter os métodos com a mesma assinatura; não usa "implements".
*/

// Animal define o comportamento esperado: quem implementar deve ter Falar() e Correr().
type Animal interface {
	Falar()
	Correr()
}

// Cachorro é um struct concreto (vazio aqui, só para o exemplo).
type Cachorro struct{}

// Falar implementa o método da interface Animal para Cachorro.
// Por isso Cachorro pode ser atribuído a uma variável do tipo Animal.
func (c Cachorro) Falar() {
	fmt.Println("Au au cachorro")
}

func (c Cachorro) Correr() {
	fmt.Println("Cachorro correndo ... ")
}

type Gato struct{}

func (g Gato) Falar() {
	fmt.Println("Gato miando ,miau miau")
}

func (g Gato) Correr() {
	fmt.Println("Gato correndo...")
}

func FalarCorrer(a Animal) {
	fmt.Print("----------- \n")
	fmt.Print("Estou em falar correr \n")
	a.Falar()
	a.Correr()

}

// Exemplo de polimorfismo: a variável é do tipo interface (Animal),
// mas em tempo de execução guarda um valor concreto (Cachorro).
func InterfacesTeste() {
	var a Animal //

	a = Cachorro{}
	b := Gato{}

	fmt.Println("Cachorro: ")

	FalarCorrer(a)

	fmt.Println("\n Gato: ")

	FalarCorrer(b)
}
