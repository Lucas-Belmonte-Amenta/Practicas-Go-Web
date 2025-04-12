package domain

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

func ProductFromProductRequest(productRequest ProductRequest, id int) Product {
	return Product{
		ID:          id,
		Name:        productRequest.Name,
		Quantity:    productRequest.Quantity,
		CodeValue:   productRequest.CodeValue,
		Expiration:  productRequest.Expiration,
		IsPublished: productRequest.IsPublished,
		Price:       productRequest.Price,
	}
}
