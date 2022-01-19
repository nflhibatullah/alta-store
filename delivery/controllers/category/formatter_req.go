package category

type CreateCategoryRequestFormat struct {
	Name string `json:"name" form:"name"`
}

type PutCategoryRequestFormat struct {
	Name string `json:"name" form:"name"`
}
