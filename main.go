package main

import (
	"altastore/configs"
	tc "altastore/delivery/controllers/transaction"
	"altastore/delivery/routes"
	tr "altastore/repository/transaction"
	"altastore/utils"

	"github.com/labstack/echo/v4"
)

func main()  {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	
	utils.InitialMigrate(db)

	transactionRepo := tr.NewTransactionRepository(db)

	transactionController := tc.NewTransactionController(transactionRepo)

	e := echo.New()

	routes.RegisterTransactionPath(e, *transactionController)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
