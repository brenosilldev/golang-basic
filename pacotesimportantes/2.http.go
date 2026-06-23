package pacotesimportantes

import (
	"fmt"
	"io"
	"net/http"
)

func HttpTest() {

	req, err := http.Get("https://www.google.com")
	if err != nil {
		panic(err)
	}

	_, err = io.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	// fmt.Println(string(res))

	fmt.Println(req.Request)

	req.Body.Close()

}
