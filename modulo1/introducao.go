package modulo1

import (
	"fmt"
	"runtime" // para pegar as informações do sistema
	"time" // para pegar a data e hora
)

// BoasVindas monta uma mensagem padrao de boas-vindas.
func BoasVindas(nome string) string {

	// Informações do sistema
	fmt.Println("=== Meu Primeiro Programa Go ===")
	fmt.Println()

	// Define valor padrão caso o nome não seja informado
	if nome == "" {
		nome = "Breno"
	}
	linguagem := "Go"
	experiencia := 3
	anos := 26

	// Printf com formatação
	fmt.Printf("👋 Olá! Eu sou %s\n", nome)
	fmt.Println("Tenhos %d nos de experiência em programação", anos)
	fmt.Printf("🔧 Estou aprendendo %s\n", linguagem + " e " + " é minha linguagem favorita")
	fmt.Printf("📅 Tenho %d anos de experiência em programação\n", experiencia)
	fmt.Println()

	// Informações do Go e do sistema
	fmt.Printf("🐹 Versão do Go: %s\n", runtime.Version())
	fmt.Printf("💻 Sistema: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("🧵 CPUs disponíveis: %d\n", runtime.NumCPU())
	fmt.Printf("⏰ Data/Hora: %s\n", time.Now().Format("02/01/2006 15:04:05"))

	
}
