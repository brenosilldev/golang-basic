package pacotesimportantes

/*
	templates = templates são usados para gerar conteúdo dinâmico a partir de dados.
	Praticamente, templates são usados para gerar HTML, mas também podem ser usados para gerar outros tipos de conteúdo, como e-mails, arquivos de configuração, etc.

*/

import (
	"os"
	"strings"
	"text/template"
)

type Curso struct {
	Nome  string
	Preco float64
}

func TemplateTest() {

	// curso := Curso{"Curso de Go", 49.99}

	// tmp := template.New("Curso Template") // template.New cria um novo template com o nome especificado

	// tmp, _ = tmp.Parse("Curso: {{.Nome}} - Preço: R$ {{.Preco}}")

	// err := tmp.Execute(os.Stdout, curso)

	// if err != nil {
	// 	panic(err)
	// }

	// // Outra forma de criar um template é usando a função template.Must, que cria um novo template e faz o parse do conteúdo em uma única linha de código.
	// t := template.Must(template.New("Curso Template").Parse("Curso: {{.Nome}} - Preço: R$ {{.Preco}}"))

	// fmt.Println("---------------------------------------------------")
	// err = t.Execute(os.Stdout, curso)

	// if err != nil {
	// 	panic(err)
	// }

	// /* --- */

	// /*
	// 	o parseFiles é usado para carregar um arquivo de template, que pode conter HTML, CSS, JS, etc. e o Execute é usado para renderizar o template com os dados fornecidos.
	// */

	// template := template.Must(template.New("index.html").ParseFiles("pacotesimportantes/templates/index.html"))

	// err = template.Execute(os.Stdout, []Curso{
	// 	{"Go", 49.99},
	// 	{"Python", 59.99},
	// 	{"JavaScript", 69.99},
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// Web Server -
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	// 	templates := []string{
	// 		"heaeder.html",
	// 		"index.html",
	// 		"footer.html",
	// 	}

	// 	template := template.Must(template.New("content.html").ParseFiles(templates...))

	// 	err := template.Execute(w, []Curso{
	// 		{"Go", 49.99},
	// 		{"Python", 59.99},
	// 		{"JavaScript", 69.99},
	// 	})

	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })

	// fmt.Println("Servidor iniciado em http://localhost:8080")
	// http.ListenAndServe(":8080", nil)

	// adicionando varios arquivos ao template, para isso é necessário criar um slice de strings com os nomes dos arquivos e passar para o ParseFiles.
	templates := []string{
		"pacotesimportantes/header.html",
		"pacotesimportantes/content.html",
		"pacotesimportantes/footer.html",
	}

	/*
		Para adicionar funções personalizadas ao template, é necessário criar um FuncMap, que é um mapa de funções, e passar para o template.New().Funcs().
	*/
	t := template.New("content.html").Funcs(template.FuncMap{"ToUpper": ToUpper})
	// ParseFiles é usado para carregar os arquivos de template, que podem conter HTML, CSS, JS, etc. e o Execute é usado para renderizar o template com os dados fornecidos.
	template := template.Must(t.ParseFiles(templates...))

	err := template.Execute(os.Stdout, []Curso{
		{"Go", 49.99},
		{"Python", 59.99},
		{"JavaScript", 69.99},
	})

	if err != nil {
		panic(err)
	}

}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}
