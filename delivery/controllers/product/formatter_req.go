package product

type CreateProductRequestFormat struct {
	Name        string `json:"name" form:"name"`
	Price       int    `json:"price" form:"price"`
	Stock       int    `json:"stock" form:"stock"`
	Description string `json:"description" form:"description"`
}

type PutProductRequestFormat struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}
