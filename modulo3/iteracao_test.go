package modulo3

import (
	"testing"
)

func TestRepetir(t *testing.T) {
	resultado := Repetir(10)
	esperado := "aaaaaaaaaa"
	if resultado != esperado {
		t.Errorf("esperado '%s', resultado '%s'", esperado, resultado)
	}
}

func BenchmarkRepetir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repetir(10)
	}
}
