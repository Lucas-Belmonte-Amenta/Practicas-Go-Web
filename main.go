package main

import (
	"fmt"
	"net/http"

	"PRACTICAS-GO-WEB/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {

	var controller *handlers.ProductController = handlers.NewProductController()

	var router chi.Router = chi.NewRouter()
	router.Get("/ping", controller.HandlerPing)
	router.Route("/products", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/", controller.HandlerGetAllProduct)
			r.Post("/", controller.HandlerCreateProduct)
			r.Get("/{id}", controller.HandlerGetProductByID)
			r.Get("/search", controller.HandlerSearchProductByPrice)
		})
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
		return
	}

}
