package repository

import (
	"PRACTICAS-GO-WEB/internal/domain"
	"PRACTICAS-GO-WEB/internal/storage"
	"cmp"
	"fmt"
	"slices"
)

type ProductRepository interface {
	GetNextID() (int, error)
	LoadAll() error
	SaveAll() error
	Get(id int) (domain.Product, error)
	GetAll() ([]domain.Product, error)
	Create(product domain.Product) (domain.Product, error)
	Update(product domain.Product) (domain.Product, error)
	Delete(id int) error
}

type productRepository struct {
	storage  storage.Storage
	products []domain.Product
	lastID   int
}

func NewProductRepository(storage storage.Storage) (*productRepository, error) {

	repository := &productRepository{storage: storage}
	err := repository.LoadAll()
	if err != nil {
		return nil, err
	}

	repository.lastID, err = repository.GetNextID()
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func (pr *productRepository) GetNextID() (int, error) {
	var product domain.Product = slices.MaxFunc(pr.products, func(p1, p2 domain.Product) int { return cmp.Compare(p1.ID, p2.ID) })
	var max int = product.ID + 1
	return max, nil
}

func (pr *productRepository) LoadAll() error {
	var products []domain.ProductStorage

	err := pr.storage.Read(&products)
	if err != nil {
		return fmt.Errorf("Error al recuperar los datos almacenados: %s", err.Error())
	}

	pr.products = domain.ProductsFromProductsStorage(products)

	return nil
}

func (pr *productRepository) SaveAll() error {

	products := domain.ProductsStorageFromProducts(pr.products)
	err := pr.storage.Write(products)
	if err != nil {
		return fmt.Errorf("Error al almacenar los datos: %s", err.Error())
	}

	return nil
}

func (pr *productRepository) Get(id int) (domain.Product, error) {

	var products, err = pr.GetAll()
	if err != nil {
		return domain.Product{}, err
	}

	index := slices.IndexFunc(products, func(p domain.Product) bool { return p.ID == id })
	if index == -1 {
		return domain.Product{}, fmt.Errorf("No se encontró el producto con el ID %d", id)
	}

	return products[index], nil

}

func (pr *productRepository) GetAll() ([]domain.Product, error) {
	return pr.products, nil
}

func (pr *productRepository) Create(product domain.Product) (domain.Product, error) {

	id, err := pr.GetNextID()
	if err != nil {
		return domain.Product{}, fmt.Errorf("Ocurrió un error durante la creación del nuevo producto: %s", err.Error())
	}
	product.ID = id

	pr.products = append(pr.products, product)
	if err := pr.SaveAll(); err != nil {
		return domain.Product{}, err
	}

	pr.lastID = id

	return product, nil
}

func (pr *productRepository) Update(product domain.Product) (domain.Product, error) {

	index := slices.IndexFunc(pr.products, func(p domain.Product) bool { return p.ID == product.ID })

	pr.products[index] = product

	if err := pr.SaveAll(); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func (pr *productRepository) Delete(id int) error {

	index := slices.IndexFunc(pr.products, func(p domain.Product) bool { return p.ID == id })
	if index == -1 {
		return fmt.Errorf("No se encontró el producto con el ID %d", id)
	}

	pr.products = slices.Delete(pr.products, index, index+1)

	if err := pr.SaveAll(); err != nil {
		return err
	}

	return nil
}
