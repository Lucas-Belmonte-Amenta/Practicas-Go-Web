package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"PRACTICAS-GO-WEB/internal/domain"
	"PRACTICAS-GO-WEB/internal/service"

	"github.com/go-chi/chi/v5"
)

type productHandler struct {
	service service.ProductService
}

type ProductHandler interface {
	HandlerPing(w http.ResponseWriter, r *http.Request)
	HandlerGetAllProduct(w http.ResponseWriter, r *http.Request)
	HandlerGetProductByID(w http.ResponseWriter, r *http.Request)
	HandlerSearchProductByPrice(w http.ResponseWriter, r *http.Request)
	HandlerCreateProduct(w http.ResponseWriter, r *http.Request)
}

// función para crear un nuevo controlador de productos
func NewProductHandler(service service.ProductService) ProductHandler {

	return &productHandler{service: service}

}

func (ph *productHandler) HandlerPing(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode("pong")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
}

func (ph *productHandler) HandlerGetAllProduct(w http.ResponseWriter, r *http.Request) {

	// Obtener todos los productos del servicio
	products, err := ph.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serializar el slice de Product a JSON y enviarlo en la respuesta
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (ph *productHandler) HandlerGetProductByID(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de los parámetros de la URL
	var idStr string = chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El ID debe ser un número entero", http.StatusBadRequest)
		return
	}

	product, err := ph.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func (ph *productHandler) HandlerSearchProductByPrice(w http.ResponseWriter, r *http.Request) {

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
	filteredProducts, err := ph.service.SearchByPrice(priceGt)

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

func (ph *productHandler) HandlerCreateProduct(w http.ResponseWriter, r *http.Request) {

	// Leer el cuerpo de la solicitud
	var productRequest domain.ProductRequest
	err := json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	var newProduct = domain.ProductFromProductRequest(productRequest)

	// Validar el producto
	productCreated, err := ph.service.Create(newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Preparar la serialización del nuevo Product a JSON y enviarlo en la respuesta despues de agregarlo al slice
	var productResponse = domain.ProductResponseFromProductBase(productCreated)
	err = json.NewEncoder(w).Encode(productResponse)
	if err != nil {
		http.Error(w, "Ocurrió un error inesperado en el procesado de la solicitud.", http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}
