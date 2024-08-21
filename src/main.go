package main

import (
	"fmt"
	"net/http"

	// "net/http"

	r "router"

	"github.com/gobuffalo/envy"
)

func main() {
	envy.Load(".env")
	DEFAULT_API_PREFIX := envy.Get("DEFAULT_API_PREFIX", "/")

	router := r.NewRouter()

	router.Use("GET", "/products", func(res http.ResponseWriter, req *r.RequestData) {
		fmt.Print("Middleware thingy!\n")
		res.Write([]byte("Middleware thingy!\n"))
	})
	
	router.Get("/products/thing", func(res http.ResponseWriter, req *r.RequestData) {
		// fmt.Fprintf(res, "Hello World!")
		res.Write([]byte("Hello World!"))
	})

	router.Get("/products/thing/{id}", func(res http.ResponseWriter, req *r.RequestData) {
		// fmt.Fprintf(res, "Hello World!")
		res.Write([]byte("Hello World!\nHere's your id: " + req.Params["id"]))
	})

	router.Listen(DEFAULT_API_PREFIX, 8080, "127.0.0.1", "Server is running")
}
