package main

import (
	"altastore/configs"
	"altastore/delivery/controllers/category"

	pc "altastore/delivery/controllers/product"
	tc "altastore/delivery/controllers/transaction"
	uc "altastore/delivery/controllers/users"
	"altastore/delivery/routes"
	cr "altastore/repository/category"
	pr "altastore/repository/product"
	tr "altastore/repository/transaction"
	ur "altastore/repository/users"
	"altastore/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	config := configs.GetConfig()

	db := utils.InitDB(config)

	utils.InitialMigrate(db)

	userRepo := ur.NewUsersRepo(db)
	transactionRepo := tr.NewTransactionRepository(db)
	productRepo := pr.NewProductRepo(db)
	categoryRepo := cr.NewCategoryRepo(db)

	userCtrl := uc.NewUsersControllers(userRepo)
	transactionController := tc.NewTransactionController(transactionRepo)
	productController := pc.NewProductControllers(productRepo)
	categoryController := category.NewCategoryControllers(categoryRepo)

	e := echo.New()

	routes.RegisterTransactionPath(e, transactionController)
	routes.RegisterUserPath(e, userCtrl)
	routes.RegisterProductPath(e, productController)
	routes.RegisterCategoryPath(e, categoryController)

	e.Logger.Fatal(e.Start(":" + config.Port))

}
