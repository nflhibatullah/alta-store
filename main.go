package main

import (
	"altastore/configs"
	"altastore/delivery/controllers/users"
	"altastore/delivery/routes"
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

	routes.RegisterPath(e, userCtrl)

	address := fmt.Sprintf("localhost:%d", config.Port)
	log.Fatal(e.Start(address))
}
