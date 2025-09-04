package main

import (
	//"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	var srv http.Server
	srv.Handler = mux
	srv.Addr = ":8080"
	srv.ListenAndServe()


	/*mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})*/
}