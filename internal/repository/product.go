package repository

import (
	"PRACTICAS-GO-WEB/internal/domain"
	"PRACTICAS-GO-WEB/pkg/utils"
	"cmp"
	"errors"
	"fmt"
	"slices"
)

type ProductRepository interface {
	GetNextID() (int, error)
	LoadAll() error
	SaveAll() error
	GetAll() ([]domain.Product, error)
	GetByID(id int) (domain.Product, error)
	Create(product domain.Product) (domain.Product, error)
}

type productRepository struct {
	products []domain.Product
	filePath string
	lastID   int
}

func NewProductRepository(filepath string) (*productRepository, error) {

	repository := &productRepository{filePath: filepath}
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
	if pr.products == nil {
		return 0, fmt.Errorf("Ocurrió un error al obtener los productos de la base de datos")
	}
	var product domain.Product = slices.MaxFunc(pr.products, func(p1, p2 domain.Product) int { return cmp.Compare(p1.ID, p2.ID) })
	var max int = product.ID + 1
	return max, nil
}

func (pr *productRepository) LoadAll() error {
	var products []domain.Product

	err := utils.ReadJSONFile(pr.filePath, &products)
	if err != nil {
		return fmt.Errorf("Error al recuperar los datos almacenados: %s", err.Error())
	}

	pr.products = products
	return nil
}

func (pr *productRepository) SaveAll() error {

	err := utils.WriteJSONFile(pr.filePath, pr.products)
	if err != nil {
		return fmt.Errorf("Error al almacenar los datos: %s", err.Error())
	}

	return nil
}

func (pr *productRepository) GetAll() ([]domain.Product, error) {
	if pr.products == nil {
		return nil, fmt.Errorf("no products found")
	}
	return pr.products, nil
}

func (pr *productRepository) GetByID(id int) (domain.Product, error) {

	index := slices.IndexFunc(pr.products, func(p domain.Product) bool { return p.ID == id })
	if index == -1 {
		return domain.Product{}, fmt.Errorf("No se encontró el producto con el ID %d", id)
	}

	return pr.products[index], nil
}

func (pr *productRepository) validateNew(newProduct domain.Product) error {

	err := newProduct.ValidateProduct()
	if err != nil {
		return err
	}

	var index int = slices.IndexFunc(pr.products, func(p domain.Product) bool {
		return (p.CodeValue == newProduct.CodeValue) || (p.ID == newProduct.ID)
	})
	if index != -1 {
		return errors.New("El código del producto ingresado ya existe")
	}

	return nil

}

func (pr *productRepository) Create(product domain.Product) (domain.Product, error) {

	if err := pr.validateNew(product); err != nil {
		return domain.Product{}, fmt.Errorf("Ocurrió un error durante la creación del nuevo producto: %s", err.Error())
	}

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
