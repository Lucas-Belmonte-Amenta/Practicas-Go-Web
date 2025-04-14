package domain

import (
	"errors"
	"time"
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

type ProductResponse struct {
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

func ProductResponseFromProductBase(product Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Quantity:    product.Quantity,
		CodeValue:   product.CodeValue,
		Expiration:  product.Expiration,
		IsPublished: product.IsPublished,
		Price:       product.Price,
	}
}

func ProductFromProductRequest(productRequest ProductRequest) Product {
	return Product{
		ID:          0,
		Name:        productRequest.Name,
		Quantity:    productRequest.Quantity,
		CodeValue:   productRequest.CodeValue,
		Expiration:  productRequest.Expiration,
		IsPublished: productRequest.IsPublished,
		Price:       productRequest.Price,
	}
}

func (product *Product) ValidateProduct() error {

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
		return errors.New("La fecha de expiración es inválida. El formato debe ser dd/mm/aaaa")
	}

	return nil

}
