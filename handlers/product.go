package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"PRACTICAS-GO-WEB/domain"
	"PRACTICAS-GO-WEB/utils"

	"github.com/go-chi/chi/v5"
)

// Controller contiene la "base de datos" en memoria
type ProductController struct {
	products []domain.Product
}

// función para crear un nuevo controlador de productos
func NewProductController() *ProductController {

	var products []domain.Product

	err := utils.ReadJSONFile("products.json", &products)
	if err != nil {
		panic("Error al leer el archivo JSON: " + err.Error())
	}

	return &ProductController{products: products}

}

func (controller *ProductController) ValidateProduct(product domain.Product) error {

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

	_, err := time.Parse("02/01/2006", product.Expiration)
	if err != nil {
		fmt.Println(err)
		return errors.New("La fecha de expiración es inválida. El formato debe ser dd/mm/aaaa")
	}

	var index int = slices.IndexFunc(controller.products, func(p domain.Product) bool { return p.CodeValue == product.CodeValue })
	if index != -1 {
		return errors.New("El código del producto ingresado ya existe")
	}

	return nil

}

func (controller *ProductController) HandlerPing(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode("pong")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
}

func (controller *ProductController) HandlerGetAllProduct(w http.ResponseWriter, r *http.Request) {

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

func (controller *ProductController) HandlerGetProductByID(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de los parámetros de la URL
	var idStr string = chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El ID debe ser un número entero", http.StatusBadRequest)
		return
	}

	// Buscar el producto en el slice usando slices.IndexFunc
	index := slices.IndexFunc(controller.products, func(product domain.Product) bool {
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

func (controller *ProductController) HandlerSearchProductByPrice(w http.ResponseWriter, r *http.Request) {

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

	filteredProducts := slices.DeleteFunc(slices.Clone(controller.products), func(product domain.Product) bool {
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

func (controller *ProductController) HandlerCreateProduct(w http.ResponseWriter, r *http.Request) {

	// Leer el cuerpo de la solicitud
	var productRequest domain.ProductRequest
	err := json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	var newProduct = domain.ProductFromProductRequest(productRequest, len(controller.products)+1)

	// Validar el producto
	err = controller.ValidateProduct(newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var productResponse = domain.ProductResponseFromProductBase(newProduct)
	// Preparar la serialización del nuevo Product a JSON y enviarlo en la respuesta despues de agregarlo al slice
	err = json.NewEncoder(w).Encode(productResponse)
	if err != nil {
		http.Error(w, "Ocurrió un error inesperado en el procesado de la solicitud.", http.StatusInternalServerError)
		return
	}

	// Agregar el producto al slice
	controller.products = append(controller.products, newProduct)

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}
