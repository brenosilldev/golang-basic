package pacotesimportantes

import (
	"net/http"
)

func FileServerTest() {
	fileserver := http.FileServer(http.Dir("./public"))
	mux := http.NewServeMux()
	mux.Handle("/", fileserver)

	http.ListenAndServe(":8080", mux)
}
