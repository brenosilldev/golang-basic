package context

import (
	"log"
	"net/http"
	"time"
)

func ContextHttpTest() {
	http.HandleFunc("/teste", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()              // Obtém o contexto da requisição HTT
	log.Println("Request Iniciado") // Log para indicar que a requisição começou

	defer log.Println("Request finalizada") // Log para indicar que a requisição terminou

	select { // select = permite aguardar múltiplas operações de comunicação (como canais) e agir com base em qual delas é concluída primeiro
	case <-time.After(3 * time.Second): // Simula um processamento que leva 5 segundos
		log.Println("Request processado com sucesso")
		w.Write([]byte("Finalizado com sucesso!")) // w.Write envia a resposta para o cliente / []byte("Finalizado com sucesso!") converte a string para um slice de bytes
	case <-ctx.Done(): // Verifica se o contexto foi cancelado (por exemplo, se o cliente cancelou a requisição)
		log.Println("Request cancelado pelo cliente:", ctx.Err())

		http.Error(w, "Request cancelado pelo cliente", http.StatusRequestTimeout) // Retorna um erro HTTP para o cliente
	}
}
