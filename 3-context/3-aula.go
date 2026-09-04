package context

import (
	"context"
	"fmt"
)

type key string

func ContextWithValueTest() {
	ctx := context.Background()

	const token key = "token"

	ctx = context.WithValue(ctx, token, 123)

	bookhotel(ctx, "João")

}

func bookhotel(ctx context.Context, name string) {
	const tokenKey key = "token"

	token := ctx.Value(tokenKey)
	if token == nil {
		fmt.Println("Token não encontrado no contexto")
		return
	}

	fmt.Printf("Reservando hotel para %s com token: %v\n", name, token)

}
