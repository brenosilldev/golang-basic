package pacotesimportantes

import (
	"fmt"
	"net/http"
)

/*
O servidor HTTP padrão do Go é muito simples de usar, mas não é muito flexível. Para criar um servidor HTTP mais flexível, podemos usar o pacote "net/http" e criar um "ServeMux" (multiplexador de solicitações).

ServeMux é um roteador de solicitações HTTP que permite registrar manipuladores de solicitações para diferentes caminhos. Ele é útil quando você deseja criar um servidor HTTP com várias rotas e manipuladores.
*/
func MuxTest() {

	mux := http.NewServeMux()

	mux.HandleFunc("/cep", BuscarCephandler)
	mux.Handle("/block", &block{title: "Ola"})

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", mux)

}

type block struct {
	title string
}

func (b *block) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(b.title))
}
