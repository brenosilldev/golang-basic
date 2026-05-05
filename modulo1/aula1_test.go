package modulo1

import (
	"testing"
)

func TestOla(t *testing.T) {

	verificaMensagemCorreta := func(t *testing.T, resultado, esperado string) {
		t.Helper() // isso é para o teste falhar e mostrar o nome da função que falhou
		if resultado != esperado {
			t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
		}
	}

	t.Run("diz olá para as pessoas em portugues", func(t *testing.T) {
		resultado := Ola("Chris", "Portugues")
		esperado := "Olá, Chris"
		verificaMensagemCorreta(t, resultado, esperado)
	})

	t.Run("'Mundo' como padrão para 'string' vazia", func(t *testing.T) {
		resultado := Ola("Mundo", "Portugues")
		esperado := "Olá, Mundo"
		verificaMensagemCorreta(t, resultado, esperado)
	})

	t.Run("diz olá para as pessoas em espanhol", func(t *testing.T) {
		resultado := Ola("Chris", "Espanhol")
		esperado := "Hola, Chris"
		verificaMensagemCorreta(t, resultado, esperado)
	})

	t.Run("diz olá para as pessoas em ingles", func(t *testing.T) {
		resultado := Ola("Chris", "Inglês")
		esperado := "Hello, Chris"
		verificaMensagemCorreta(t, resultado, esperado)
	})
}
