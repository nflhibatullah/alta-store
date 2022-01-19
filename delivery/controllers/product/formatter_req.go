package product

type CreateProductRequestFormat struct {
	Name        string `json:"name" form:"name"`
	Price       int    `json:"price" form:"price"`
	Stock       int    `json:"stock" form:"stock"`
	Description string `json:"description" form:"description"`
	CategoryID  int    `json:"categoryID" form:"categoryID"`
}

type PutProductRequestFormat struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}
