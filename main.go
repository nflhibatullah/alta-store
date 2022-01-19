package main

import (
	"altastore/configs"
	tc "altastore/delivery/controllers/transaction"
	uc "altastore/delivery/controllers/users"
	"altastore/delivery/routes"
	tr "altastore/repository/transaction"
	ur "altastore/repository/users"
	"altastore/utils"

	"github.com/labstack/echo/v4"
)

func main()  {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	
	utils.InitialMigrate(db)

	userRepo := ur.NewUsersRepo(db)
	transactionRepo := tr.NewTransactionRepository(db)

	userCtrl := uc.NewUsersControllers(userRepo)
	transactionController := tc.NewTransactionController(transactionRepo)

	e := echo.New()

	routes.RegisterTransactionPath(e, transactionController)
	routes.RegisterUserPath(e, userCtrl)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
