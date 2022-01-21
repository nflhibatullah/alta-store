package category

type CreateCategoryRequestFormat struct {
	Name string `json:"name" form:"name"  validate:"required"`
}

type PutCategoryRequestFormat struct {
	Name string `json:"name" form:"name" validate:"required"`
}
