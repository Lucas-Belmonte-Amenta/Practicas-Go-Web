package service

import (
	"PRACTICAS-GO-WEB/internal/domain"
	"PRACTICAS-GO-WEB/internal/repository"

	"errors"
	"fmt"
	"slices"
)

type ProductService interface {
	GetProducts() ([]domain.ProductResponse, error)
	GetProductByID(id int) (domain.ProductResponse, error)
	SearchProductByPrice(priceGt float64) ([]domain.ProductResponse, error)
	PostProduct(product domain.ProductRequest) (domain.ProductResponse, error)
	PutProduct(id int, product domain.ProductRequest) (domain.ProductResponse, error)
	PatchProduct(id int, product domain.ProductRequest) (domain.ProductResponse, error)
	DeleteProduct(id int) error
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

func (ps *productService) GetProducts() ([]domain.ProductResponse, error) {

	var products, err = ps.productRepository.GetAll()

	if err != nil {
		return nil, err
	}

	productsResponse := domain.ProductResponsesFromProductsBase(products)

	return productsResponse, nil

}

func (ps *productService) GetProductByID(id int) (domain.ProductResponse, error) {

	product, err := ps.productRepository.Get(id)
	if err != nil {
		return domain.ProductResponse{}, err
	}

	productResponse := domain.ProductResponseFromProductBase(product)

	return productResponse, nil

}

func (ps *productService) SearchProductByPrice(priceGt float64) ([]domain.ProductResponse, error) {

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

	productsResponses := domain.ProductResponsesFromProductsBase(filteredProducts)

	return productsResponses, nil

}

func (ps *productService) validateCodeValue(codeValue string) error {

	var products, err = ps.productRepository.GetAll()
	if err != nil {
		return err
	}

	var index int = slices.IndexFunc(products, func(product domain.Product) bool {
		return product.CodeValue == codeValue
	})

	if index == -1 {
		return nil
	}

	return fmt.Errorf("Ya existe un producto registrado con el codigo %s", codeValue)
}

func (ps *productService) validateNewProduct(newProduct domain.Product) error {

	err := newProduct.ValidateProduct()
	if err != nil {
		return err
	}

	err = ps.validateCodeValue(newProduct.CodeValue)
	if err != nil {
		return err
	}

	return nil

}

func (ps *productService) PostProduct(product domain.ProductRequest) (domain.ProductResponse, error) {

	newProduct := domain.ProductFromProductRequest(product)

	id, err := ps.productRepository.GetNextID()
	if err != nil {
		return domain.ProductResponse{}, fmt.Errorf("Error al crear un nuevo producto: %s", err.Error())
	}
	newProduct.ID = id

	if err := ps.validateNewProduct(newProduct); err != nil {
		return domain.ProductResponse{}, fmt.Errorf("Ocurrió un error durante la creación del nuevo producto: %s", err.Error())
	}

	productCreated, err := ps.productRepository.Create(newProduct)
	if err != nil {
		return domain.ProductResponse{}, fmt.Errorf("Error al crear un nuevo producto: %s", err.Error())
	}

	return domain.ProductResponseFromProductBase(productCreated), nil

}

func (ps *productService) PutProduct(id int, product domain.ProductRequest) (domain.ProductResponse, error) {

	productToUpdate := domain.ProductFromProductRequest(product)

	if _, err := ps.GetProductByID(id); err != nil {
		return domain.ProductResponse{}, err
	}

	productToUpdate.ID = id

	if err := productToUpdate.ValidateProduct(); err != nil {
		return domain.ProductResponse{}, err
	}

	productUpdated, err := ps.productRepository.Update(productToUpdate)
	if err != nil {
		return domain.ProductResponse{}, err
	}

	return domain.ProductResponseFromProductBase(productUpdated), nil
}

func (ps *productService) PatchProduct(id int, product domain.ProductRequest) (domain.ProductResponse, error) {

	oldProduct, err := ps.productRepository.Get(id)
	if err != nil {
		return domain.ProductResponse{}, err
	}

	if product.Name != nil {
		oldProduct.Name = *product.Name
	}

	if product.Quantity != nil {
		oldProduct.Quantity = *product.Quantity
	}

	if (product.CodeValue != nil) && (*product.CodeValue != oldProduct.CodeValue) {
		if ps.validateCodeValue(*product.CodeValue) != nil {
			oldProduct.CodeValue = *product.CodeValue
		}
	}

	if product.Expiration != nil {
		oldProduct.Expiration = *product.Expiration
	}

	if product.Price != nil {
		oldProduct.Price = *product.Price
	}

	if oldProduct.ValidateProduct() != nil {
		return domain.ProductResponse{}, err
	}

	productUpdated, err := ps.productRepository.Update(oldProduct)
	if err != nil {
		return domain.ProductResponse{}, err
	}

	return domain.ProductResponseFromProductBase(productUpdated), nil

}

func (ps *productService) DeleteProduct(id int) error {

	if _, err := ps.GetProductByID(id); err != nil {
		return err
	}

	err := ps.productRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil

}
