package fundacao

import "fmt"

type Conta struct {
	saldo          int
	saldodepositar int
	saldoaretirar  int
}

func (c *Conta) Depositar(valor int) int {
	c.saldo += valor
	println("Saldo atual depositado: ", c.saldo)
	return c.saldo

}

func (c *Conta) Retirar(valor int) int {
	c.saldo -= valor
	println("Saldo atual retirado: ", c.saldo)
	return c.saldo
}

func PonteirosStructsTests() {
	conta := Conta{
		saldo: 300,
	}

	conta.saldodepositar = 100

	conta.Depositar(conta.saldodepositar)
	fmt.Printf("Saldo atual: %d\n", conta.saldo)

	conta.saldoaretirar = 10
	conta.Retirar(conta.saldoaretirar)

	fmt.Printf("Saldo atual: %d\n", conta.saldo)

}
