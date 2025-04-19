package domain

import (
	"errors"
	"fmt"
	"time"
)

type Product struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Quantity    int        `json:"quantity"`
	CodeValue   string     `json:"code_value"`
	Expiration  *time.Time `json:"expiration_date,omitempty"`
	IsPublished bool       `json:"is_published,omitempty"`
	Price       float64    `json:"price"`
}

type ProductStorage struct {
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
	Expiration  *string `json:"expiration_date,omitempty"`
	IsPublished bool    `json:"is_published,omitempty"`
	Price       float64 `json:"price"`
}

type ProductRequest struct {
	Name        *string  `json:"name"`
	Quantity    *int     `json:"quantity"`
	CodeValue   *string  `json:"code_value"`
	Expiration  *string  `json:"expiration_date,omitempty"`
	IsPublished *bool    `json:"is_published,omitempty"`
	Price       *float64 `json:"price"`
}

func ProductResponseFromProductBase(product Product) ProductResponse {

	var expiration *string
	if product.Expiration != nil {
		timeStr := product.Expiration.Format("02/01/2006")
		expiration = &timeStr
	} else {
		expiration = nil
	}

	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Quantity:    product.Quantity,
		CodeValue:   product.CodeValue,
		Expiration:  expiration,
		IsPublished: product.IsPublished,
		Price:       product.Price,
	}

}

func ProductResponsesFromProductsBase(products []Product) []ProductResponse {
	productsResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productsResponses[i] = ProductResponseFromProductBase(product)
	}
	return productsResponses
}

func ProductFromProductRequest(productRequest ProductRequest) (Product, error) {

	if productRequest.Name == nil {
		defaultName := ""
		productRequest.Name = &defaultName
	}

	if productRequest.Quantity == nil {
		defaultQuantity := 0
		productRequest.Quantity = &defaultQuantity
	}

	if productRequest.CodeValue == nil {
		defaultCodeValue := ""
		productRequest.CodeValue = &defaultCodeValue
	}

	var expiration *time.Time
	if productRequest.Expiration != nil {
		timeStr, err := time.Parse("02/01/2006", *productRequest.Expiration)
		if err != nil {
			return Product{}, fmt.Errorf("Error al parsear la fecha de expiración: %s", err.Error())
		}
		expiration = &timeStr
	} else {
		expiration = nil
	}

	if productRequest.IsPublished == nil {
		defaultIsPublished := false
		productRequest.IsPublished = &defaultIsPublished
	}

	if productRequest.Price == nil {
		defaultPrice := 0.0
		productRequest.Price = &defaultPrice
	}

	product := Product{
		ID:          0,
		Name:        *productRequest.Name,
		Quantity:    *productRequest.Quantity,
		CodeValue:   *productRequest.CodeValue,
		Expiration:  expiration,
		IsPublished: *productRequest.IsPublished,
		Price:       *productRequest.Price,
	}

	return product, nil

}

func (product *Product) ValidateProduct() error {

	switch {
	case product.Name == "":
		return errors.New("El nombre del producto es un campo requerido")
	case product.CodeValue == "":
		return errors.New("El codigo del producto es un campo requerido")
	case product.Price == 0:
		return errors.New("El precio del producto es un campo requerido")
	case product.Quantity == 0:
		return errors.New("La stock del producto es un campo requerido")
	}

	return nil

}

func ProductsFromProductsStorage(productsStorage []ProductStorage) []Product {

	var products []Product

	for _, productStorage := range productsStorage {

		var expiration *time.Time
		if productStorage.Expiration != "" {
			timeStr, _ := time.Parse("02/01/2006", productStorage.Expiration)
			expiration = &timeStr
		} else {
			expiration = nil
		}

		product := Product{
			ID:          productStorage.ID,
			Name:        productStorage.Name,
			Quantity:    productStorage.Quantity,
			CodeValue:   productStorage.CodeValue,
			Expiration:  expiration,
			IsPublished: productStorage.IsPublished,
			Price:       productStorage.Price,
		}

		products = append(products, product)

	}

	return products

}

func ProductsStorageFromProducts(products []Product) []ProductStorage {

	var productsStorage []ProductStorage

	for _, product := range products {

		var expiration string
		if product.Expiration != nil {
			timeStr := product.Expiration.Format("02/01/2006")
			expiration = timeStr
		} else {
			expiration = ""
		}

		productStorage := ProductStorage{
			ID:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			Expiration:  expiration,
			IsPublished: product.IsPublished,
			Price:       product.Price,
		}

		productsStorage = append(productsStorage, productStorage)

	}

	return productsStorage

}
