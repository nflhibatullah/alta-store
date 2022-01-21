package main

import (
	"altastore/configs"
	"altastore/delivery/common"
	ctc "altastore/delivery/controllers/cart"
	cc "altastore/delivery/controllers/category"
	"altastore/delivery/middlewares"

	"github.com/go-playground/validator/v10"

	pc "altastore/delivery/controllers/product"
	tc "altastore/delivery/controllers/transaction"
	uc "altastore/delivery/controllers/users"
	"altastore/delivery/routes"
	ctr "altastore/repository/cart"
	cr "altastore/repository/category"
	pr "altastore/repository/product"
	tr "altastore/repository/transaction"
	ur "altastore/repository/users"
	"altastore/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := configs.GetConfig()

	db := utils.InitDB(config)

	utils.InitialMigrate(db)

	userRepo := ur.NewUsersRepo(db)
	transactionRepo := tr.NewTransactionRepository(db)
	productRepo := pr.NewProductRepo(db)
	categoryRepo := cr.NewCategoryRepo(db)
	cartRepo := ctr.NewCartRepository(db)

	userCtrl := uc.NewUsersControllers(userRepo)
	transactionController := tc.NewTransactionController(transactionRepo)
	productController := pc.NewProductControllers(productRepo)
	categoryController := cc.NewCategoryControllers(categoryRepo)
	cartController := ctc.NewCartController(cartRepo)

	e := echo.New()

	middlewares.LogMiddleware(e)

	e.Pre(middleware.RemoveTrailingSlash())

	e.Validator = &common.CustomValidator{Validator: validator.New()}

	routes.RegisterTransactionPath(e, transactionController)
	routes.RegisterUserPath(e, userCtrl)
	routes.RegisterProductPath(e, productController)
	routes.RegisterCategoryPath(e, categoryController)
	routes.RegisterCartPath(e, cartController)

	e.Logger.Fatal(e.Start(":" + config.Port))

}
