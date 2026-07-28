package pacotesimportantes

import (
	"os"
	"text/template"
)

type Curso struct {
	Nome  string
	Preco float64
}

func TemplateTest() {
	curso := Curso{"Curso de Go", 49.99}
	tmp := template.New("Curso Template")
	tmp, _ = tmp.Parse("Curso: {{.Nome}} - Preço: R$ {{.Preco}}")

	err := tmp.Execute(os.Stdout, curso)

	if err != nil {
		panic(err)
	}
}
