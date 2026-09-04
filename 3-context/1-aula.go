package context

import (
	"context"
	"log"
	"time"
)

// O context é usado para controlar operações assíncronas, como cancelamento e timeout. Ele é útil para gerenciar o ciclo de vida de goroutines e operações que podem levar tempo.

func AulaTeste() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel() // O cancelamento é chamado para liberar recursos quando a função termina

	bookHotel(ctx)
}

func bookHotel(ctx context.Context) {
	// Simula uma operação de reserva de hotel que leva algum tempo
	select {
	case <-ctx.Done(): // Verifica se o contexto foi cancelado ou expirou
		log.Println("Booking canceled:", ctx.Err())
	case <-time.After(3 * time.Second):
		// Proceed with booking logic
		log.Println("Hotel booked successfully")
	}
}
