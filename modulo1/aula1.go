package modulo1

const prefixoOlaPortugues = "Olá, "

const (
	Portugues = "Portugues"
	Espanhol  = "Espanhol"
	Inglês    = "Inglês"
	Francês   = "Francês"
	Japonês   = "Japonês"
	Coreano   = "Coreano"
	Mandarim  = "Mandarim"
	Arabe     = "Arabe"
	Hebreu    = "Hebreu"
	Russo     = "Russo"
)

func Ola(nome string, idioma string) string {
	if idioma == Portugues {
		return "Olá, " + nome
	}
	if idioma == Espanhol {
		return "Hola, " + nome
	}
	if idioma == Inglês {
		return "Hello, " + nome
	}
	if idioma == Francês {
		return "Bonjour, " + nome
	}
	if idioma == Japonês {
		return "こんにちは, " + nome
	}
	if idioma == Coreano {
		return "안녕하세요, " + nome
	}
	if idioma == Mandarim {
		return "你好, " + nome
	}
	if idioma == Arabe {
		return "مرحبا, " + nome
	}
	if idioma == Hebreu {
		return "שלום, " + nome
	}
	if idioma == Russo {
		return "Здравствуйте, " + nome
	}
	return "Olá, " + nome
}
