package pacotesimportantes

import (
	"fmt"
	"net/http"
	"os"
)

func BuscaCepTest() {

	for _, url := range os.Args[1:] {

		req, err := http.Get(url)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n ", err)
		}

		defer req.Body.Close()

		res, err :+ req.

		println(url)
	}
}
