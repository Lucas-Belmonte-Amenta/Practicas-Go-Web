package service

import (
	"PRACTICAS-GO-WEB/internal/domain"
	"PRACTICAS-GO-WEB/internal/repository"
	"errors"
	"fmt"
	"slices"
)

type ProductService interface {
	GetAll() ([]domain.Product, error)
	GetByID(id int) (domain.Product, error)
	SearchByPrice(priceGt float64) ([]domain.Product, error)
	Create(product domain.Product) (domain.Product, error)
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) (*productService, error) {

	if productRepository == nil {
		return nil, errors.New("productRepository is required")
	}

	return &productService{productRepository: productRepository}, nil

}

func (ps *productService) GetAll() ([]domain.Product, error) {
	var products, err = ps.productRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps *productService) GetByID(id int) (domain.Product, error) {
	var product, err = ps.productRepository.GetByID(id)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (ps *productService) SearchByPrice(priceGt float64) ([]domain.Product, error) {

	var products, err = ps.productRepository.GetAll()
	if err != nil {
		return nil, err
	}

	filteredProducts := slices.DeleteFunc(slices.Clone(products), func(product domain.Product) bool {
		return !(product.Price >= priceGt)
	})
	if len(filteredProducts) == 0 {
		return nil, fmt.Errorf("No se encontraron productos con un precio mayor o igual a %f", priceGt)
	}

	return filteredProducts, nil

}

func (ps *productService) Create(product domain.Product) (domain.Product, error) {

	id, err := ps.productRepository.GetNextID()
	if err != nil {
		return domain.Product{}, fmt.Errorf("Error al crear un nuevo producto: %s", err.Error())
	}
	product.ID = id

	return ps.productRepository.Create(product)

}
