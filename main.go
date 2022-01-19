package main

import (
	"altastore/configs"
	"altastore/delivery/controllers/category"
	"altastore/delivery/controllers/product"
	"altastore/delivery/controllers/users"
	"altastore/delivery/routes"
	categoryRepo "altastore/repository/category"
	productRepo "altastore/repository/product"
	userRepo "altastore/repository/users"
	"altastore/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	config := configs.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	userRepo := userRepo.NewUsersRepo(db)
	userCtrl := user.NewUsersControllers(userRepo)

	productRepo := productRepo.NewProductRepo(db)
	productCtrl := product.NewProductControllers(productRepo)

	categoryRepo := categoryRepo.NewCategoryRepo(db)
	categoryCtrl := category.NewCategoryControllers(categoryRepo)

	routes.RegisterPath(e, userCtrl, productCtrl, categoryCtrl)

	address := fmt.Sprintf("localhost:%d", config.Port)
	log.Fatal(e.Start(address))
}
