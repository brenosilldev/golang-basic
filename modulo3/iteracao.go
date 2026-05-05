package modulo3

func Repetir(n int) string {
	var repeticoes string
	for i := 0; i < n; i++ {
		repeticoes += "a"
	}
	return repeticoes
}
