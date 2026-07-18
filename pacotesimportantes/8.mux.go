package pacotesimportantes

import (
	"fmt"
	"net/http"
)

// ServerMux é uma estrutura que permite criar um roteador HTTP personalizado em Go. Ele permite registrar manipuladores de rotas para diferentes caminhos e métodos HTTP, facilitando a criação de APIs e servidores web. O ServeMux é útil para organizar o código e gerenciar diferentes endpoints de forma eficiente.
func MuxTest() {

	mux := http.NewServeMux()

	mux.HandleFunc("/cep", BuscarCephandler)
	fmt.Println("http://localhost:8080")

	http.ListenAndServe(":8080", mux)

}
