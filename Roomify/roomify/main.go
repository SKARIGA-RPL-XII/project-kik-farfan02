package main

import (
	"log"
	
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/config"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/handler"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/repository"
	router "github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/route"
	service "github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()
	config.ConnectDB(cfg)

	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(*userRepo, cfg)
	userHandler := handler.NewUserHandler(*userService)

	deptRepo := repository.NewDeptRepository(config.DB)
	deptService := service.NewDeptService(*deptRepo)
	deptHandler := handler.NewDeptHandler(*deptService)

	authRepo := repository.NewAuthRepository(config.DB)
	authService := service.NewAuthService(*authRepo, cfg)
	authHandler := handler.NewAuthHandler(*authService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New())

	router.Router(app, userHandler, deptHandler, authHandler)

	log.Fatal(app.Listen("0.0.0.0:3000"))
}