package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	Expiration  string  `json:"expiration_date"`
	IsPublished bool    `json:"is_published"`
	Price       float64 `json:"price"`
}

var Products []Product = []Product{}

func ReadJSONFile[T any](fileName string, emptyListEntity *[]T) error {
	// Abrir el archivo
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Error al abrir el archivo Json: %v\n", err)
	}
	defer file.Close()

	// Leer el contenido del archivo
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Error al leer el archivo Json: %v\n", err)
	}

	// Deserializar el JSON en el slice de Product
	err = json.Unmarshal(byteValue, emptyListEntity)
	if err != nil {
		return fmt.Errorf("Error al deserializar el JSON: %v\n", err)
	}

	return nil

}

func HandlerPing(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode("pong")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
}

func HandlerGetAllProduct(w http.ResponseWriter, r *http.Request) {

	// Serializar el slice de Product a JSON y enviarlo en la respuesta
	err := json.NewEncoder(w).Encode(Products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func HandlerGetProductByID(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de los parámetros de la URL
	var idStr string = chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El ID debe ser un número entero", http.StatusBadRequest)
		return
	}

	// Buscar el producto en el slice usando slices.IndexFunc
	index := slices.IndexFunc(Products, func(product Product) bool {
		return product.ID == id
	})

	// Obtener el producto encontrado o devolver un error si no se encontró
	if index == -1 {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}
	product := Products[index]

	// Serializar el slice de Product a JSON y enviarlo en la respuesta
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func HandlerSearchProductByPrice(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de los parámetros de la URL
	var priceGtStr string = r.URL.Query().Get("priceGt")
	if priceGtStr == "" {
		http.Error(w, "El priceGt es requerido", http.StatusBadRequest)
		return
	}

	// Convertir el priceGt a float64
	priceGt, err := strconv.ParseFloat(priceGtStr, 64)
	if err != nil {
		http.Error(w, "El priceGt debe ser un numero decimal", http.StatusBadRequest)
		return
	}

	// Buscar los productos en el slice

	filteredProducts := slices.DeleteFunc(slices.Clone(Products), func(product Product) bool {
		return !(product.Price >= priceGt)
	})

	if len(filteredProducts) < 1 {
		http.Error(w, "No se encontraron productos", http.StatusNotFound)
		return
	}

	// Serializar el slice de Product a JSON y enviarlo en la respuesta
	err = json.NewEncoder(w).Encode(filteredProducts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func main() {

	err := ReadJSONFile("products.json", &Products)
	if err != nil {
		panic("Error al leer el archivo JSON: " + err.Error())
	}

	var router chi.Router = chi.NewRouter()
	router.Get("/ping", HandlerPing)
	router.Route("/products", func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Get("/", HandlerGetAllProduct)
			r.Get("/{id}", HandlerGetProductByID)
			r.Get("/search", HandlerSearchProductByPrice)

		})

	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
		return
	}

}
