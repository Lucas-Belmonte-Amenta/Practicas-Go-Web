package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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
	HandlerUpdateProduct(w http.ResponseWriter, r *http.Request)
	HandlerUpdatePartialProduct(w http.ResponseWriter, r *http.Request)
	HandlerDeleteProduct(w http.ResponseWriter, r *http.Request)
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
	products, err := ph.service.GetProducts()
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

	product, err := ph.service.GetProductByID(id)
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
	filteredProducts, err := ph.service.SearchProductByPrice(priceGt)

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

	// Validar los campos del producto de la solicitud
	err = ph.validateFullRequest(productRequest)
	if err != nil {
		errStr := fmt.Sprintf("Error al validar los datos del producto: %s", err.Error())
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	// Validar el producto
	productCreated, err := ph.service.PostProduct(productRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Preparar la serialización del nuevo Product a JSON y enviarlo en la respuesta despues de agregarlo al slice
	err = json.NewEncoder(w).Encode(productCreated)
	if err != nil {
		http.Error(w, "Ocurrió un error inesperado en el procesado de la solicitud.", http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func (ph *productHandler) HandlerUpdateProduct(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de los parámetros de la URL
	var idStr string = chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El ID debe ser un número entero", http.StatusBadRequest)
		return
	}

	// Leer el cuerpo de la solicitud
	var productRequest domain.ProductRequest
	err = json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Validar los campos del producto de la solicitud
	err = ph.validateFullRequest(productRequest)
	if err != nil {
		errStr := fmt.Sprintf("Error al validar los datos del producto: %s", err.Error())
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	// Validar el producto
	productUpdated, err := ph.service.PutProduct(id, productRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Preparar la serialización del nuevo Product a JSON y enviarlo en la respuesta despues de agregarlo al slice
	err = json.NewEncoder(w).Encode(productUpdated)
	if err != nil {
		http.Error(w, "Ocurrió un error inesperado en el procesado de la solicitud.", http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func (ph *productHandler) HandlerUpdatePartialProduct(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de los parámetros de la URL
	var idStr string = chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El ID debe ser un número entero", http.StatusBadRequest)
		return
	}

	// Leer el cuerpo de la solicitud
	var productRequest domain.ProductRequest
	err = json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Validar el producto
	productUpdated, err := ph.service.PatchProduct(id, productRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Preparar la serialización del nuevo Product a JSON y enviarlo en la respuesta despues de agregarlo al slice
	err = json.NewEncoder(w).Encode(productUpdated)
	if err != nil {
		http.Error(w, "Ocurrió un error inesperado en el procesado de la solicitud.", http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func (ph *productHandler) HandlerDeleteProduct(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de los parámetros de la URL
	var idStr string = chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El ID debe ser un número entero", http.StatusBadRequest)
		return
	}

	err = ph.service.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode("El producto se eliminó correctamente")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado Content-Type y el código de estado
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

}

func (ph *productHandler) validateFullRequest(productRequest domain.ProductRequest) error {

	if productRequest.Name == nil {
		return errors.New("El nombre del producto es un campo requerido")
	}

	if productRequest.Quantity == nil {
		return errors.New("El stock del producto es un campo requerido")
	}

	if productRequest.CodeValue == nil {
		return errors.New("El codigo del producto es un campo requerido")
	}

	if productRequest.Expiration == nil {
		return errors.New("La fecha de expiración del producto es un campo requerido")
	}

	if productRequest.IsPublished == nil {
		return errors.New("El estado de publicación del producto es un campo requerido")
	}

	if productRequest.Price == nil {
		return errors.New("El precio del producto es un campo requerido")
	}

	return nil

}
