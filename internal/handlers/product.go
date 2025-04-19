package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"PRACTICAS-GO-WEB/internal/domain"
	"PRACTICAS-GO-WEB/internal/service"
	"PRACTICAS-GO-WEB/pkg/web"

	"github.com/go-chi/chi/v5"
)

type productHandler struct {
	service            service.ProductService
	tokenAuthorization string
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

	token := os.Getenv("Token")
	return &productHandler{service: service, tokenAuthorization: token}

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
		web.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.Success(w, http.StatusOK, "products found", products)

}

func (ph *productHandler) HandlerGetProductByID(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de los parámetros de la URL
	id, err := ph.validateHeaderID(w, r)
	if err != nil {
		return
	}

	product, err := ph.service.GetProductByID(id)
	if err != nil {
		web.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	web.Success(w, http.StatusOK, "product found", product)

}

func (ph *productHandler) HandlerSearchProductByPrice(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de los parámetros de la URL
	var priceGtStr string = r.URL.Query().Get("priceGt")
	if priceGtStr == "" {
		web.Error(w, http.StatusBadRequest, "El valor de priceGt es requerido")
		return
	}

	// Convertir el priceGt a float64
	priceGt, err := strconv.ParseFloat(priceGtStr, 64)
	if err != nil {
		web.Error(w, http.StatusBadRequest, "El valor de priceGt debe ser un numero decimal")
		return
	}

	// Buscar los productos en el slice
	filteredProducts, err := ph.service.SearchProductByPrice(priceGt)

	web.Success(w, http.StatusOK, "products found", filteredProducts)

}

func (ph *productHandler) HandlerCreateProduct(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Token")
	if token != ph.tokenAuthorization {
		web.Error(w, http.StatusUnauthorized, "Token de autentificación inválido")
		return
	}

	// Leer el cuerpo de la solicitud
	var productRequest domain.ProductRequest
	err := json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		web.Error(w, http.StatusBadRequest, "Error al leer el cuerpo de la solicitud")
		return
	}

	// Validar los campos del producto de la solicitud
	err = ph.validateFullRequest(productRequest)
	if err != nil {
		errStr := fmt.Sprintf("Error al validar los datos del producto: %s", err.Error())
		web.Error(w, http.StatusBadRequest, errStr)
		return
	}

	// Validar el producto
	productCreated, err := ph.service.PostProduct(productRequest)
	if err != nil {
		errStr := fmt.Sprintf("Error al registrar el nuevo producto: %s", err.Error())
		web.Error(w, http.StatusInternalServerError, errStr)
		return
	}

	web.Success(w, http.StatusCreated, "product created", productCreated)

}

func (ph *productHandler) HandlerUpdateProduct(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Token")
	if token != ph.tokenAuthorization {
		web.Error(w, http.StatusUnauthorized, "Token de autentificación inválido")
		return
	}

	// Obtener el ID de los parámetros de la URL
	id, err := ph.validateHeaderID(w, r)
	if err != nil {
		return
	}

	// Leer el cuerpo de la solicitud
	var productRequest domain.ProductRequest
	err = json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		errStr := fmt.Sprintf("Error al leer el cuerpo de la solicitud: %s", err.Error())
		web.Error(w, http.StatusBadRequest, errStr)
		return
	}

	// Validar los campos del producto de la solicitud
	err = ph.validateFullRequest(productRequest)
	if err != nil {
		errStr := fmt.Sprintf("Error al validar los datos del producto: %s", err.Error())
		web.Error(w, http.StatusBadRequest, errStr)
		return
	}

	// Validar el producto
	productUpdated, err := ph.service.PutProduct(id, productRequest)
	if err != nil {
		errStr := fmt.Sprintf("Error al actualizar el producto: %s", err.Error())
		web.Error(w, http.StatusBadRequest, errStr)
		return
	}

	web.Success(w, http.StatusOK, "product updated", productUpdated)

}

func (ph *productHandler) HandlerUpdatePartialProduct(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Token")
	if token != ph.tokenAuthorization {
		web.Error(w, http.StatusUnauthorized, "Token de autentificación inválido")
		return
	}

	// Obtener el ID de los parámetros de la URL
	id, err := ph.validateHeaderID(w, r)
	if err != nil {
		return
	}

	// Leer el cuerpo de la solicitud
	var productRequest domain.ProductRequest
	err = json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		web.Error(w, http.StatusBadRequest, "Error al leer el cuerpo de la solicitud")
		return
	}

	// Validar el producto
	productUpdated, err := ph.service.PatchProduct(id, productRequest)
	if err != nil {
		errStr := fmt.Sprintf("Error al actualizar el producto: %s", err.Error())
		web.Error(w, http.StatusBadRequest, errStr)
		return
	}

	web.Success(w, http.StatusOK, "product updated", productUpdated)

}

func (ph *productHandler) HandlerDeleteProduct(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Token")
	if token != ph.tokenAuthorization {
		web.Error(w, http.StatusUnauthorized, "Token de autentificación inválido")
		return
	}

	// Obtener el ID de los parámetros de la URL
	id, err := ph.validateHeaderID(w, r)
	if err != nil {
		return
	}

	err = ph.service.DeleteProduct(id)
	if err != nil {
		errStr := fmt.Sprintf("Error al eliminar el producto: %s", err.Error())
		web.Error(w, http.StatusInternalServerError, errStr)
		return
	}

	web.Success(w, http.StatusNoContent, "product deleted", nil)

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

	if productRequest.Expiration != nil {
		if _, err := time.Parse("02/01/2006", *productRequest.Expiration); err != nil {
			return errors.New("La fecha de expiración no posee un formato válido")
		}
	}

	if productRequest.IsPublished == nil {
		return errors.New("El estado de publicación del producto es un campo requerido")
	}

	if productRequest.Price == nil {
		return errors.New("El precio del producto es un campo requerido")
	}

	return nil

}

func (ph *productHandler) validateHeaderID(w http.ResponseWriter, r *http.Request) (int, error) {

	// Obtener el ID de los parámetros de la URL
	var idStr string = chi.URLParam(r, "id")
	if idStr == "" {
		web.Error(w, http.StatusBadRequest, "El ID del producto es un campo requerido")
		return 0, errors.New("El ID es un campo requerido")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		web.Error(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return 0, err
	}

	return id, nil

}
