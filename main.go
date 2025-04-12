package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	Expiration  string  `json:"expiration_date"`
	IsPublished bool    `json:"is_published,omitempty"`
	Price       float64 `json:"price"`
}

type ProductRequest struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	Expiration  string  `json:"expiration_date"`
	IsPublished bool    `json:"is_published,omitempty"`
	Price       float64 `json:"price"`
}

// Controller contiene la "base de datos" en memoria
type productController struct {
	products []Product
}

// función para leer un archivo JSON y deserializarlo en un slice de Product
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

// función para crear un nuevo controlador de productos
func NewProductController() *productController {

	var products []Product

	err := ReadJSONFile("products.json", &products)
	if err != nil {
		panic("Error al leer el archivo JSON: " + err.Error())
	}

	return &productController{products: products}

}

func (controller *productController) ValidateProduct(product Product) error {

	switch {
	case product.Name == "":
		return errors.New("El nombre del producto es un campo requerido")
	case product.CodeValue == "":
		return errors.New("El codigo del producto es un campo requerido")
	case product.Expiration == "":
		return errors.New("La fecha de expiración del producto es un campo requerido")
	case product.Price == 0:
		return errors.New("El precio del producto es un campo requerido")
	case product.Quantity == 0:
		return errors.New("La stock del producto es un campo requerido")
	}

	_, err := time.Parse("25-05-1810", product.Expiration)
	if err != nil {
		return errors.New("La fecha de expiración es inválida. El formato debe ser dd-mm-aaaa")
	}

	var index int = slices.IndexFunc(controller.products, func(p Product) bool { return p.CodeValue == product.CodeValue })
	if index != -1 {
		return errors.New("El código del producto ingresado ya existe")
	}

	return nil

}

func (controller *productController) HandlerPing(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode("pong")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
}

func (controller *productController) HandlerGetAllProduct(w http.ResponseWriter, r *http.Request) {

	// Serializar el slice de Product a JSON y enviarlo en la respuesta
	err := json.NewEncoder(w).Encode(controller.products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (controller *productController) HandlerGetProductByID(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de los parámetros de la URL
	var idStr string = chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El ID debe ser un número entero", http.StatusBadRequest)
		return
	}

	// Buscar el producto en el slice usando slices.IndexFunc
	index := slices.IndexFunc(controller.products, func(product Product) bool {
		return product.ID == id
	})

	// Obtener el producto encontrado o devolver un error si no se encontró
	if index == -1 {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}
	product := controller.products[index]

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

func (controller *productController) HandlerSearchProductByPrice(w http.ResponseWriter, r *http.Request) {

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

	filteredProducts := slices.DeleteFunc(slices.Clone(controller.products), func(product Product) bool {
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

	var controller *productController = NewProductController()

	var router chi.Router = chi.NewRouter()
	router.Get("/ping", controller.HandlerPing)
	router.Route("/products", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/", controller.HandlerGetAllProduct)
			r.Get("/{id}", controller.HandlerGetProductByID)
			r.Get("/search", controller.HandlerSearchProductByPrice)
		})
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
		return
	}

}
