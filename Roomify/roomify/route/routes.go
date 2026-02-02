package router

import (
	"log"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/config"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/handler"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/middleware"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App, userhandler *handler.UserHandler, depthandler *handler.DeptHandler) {
	// Load config sekali untuk middleware
	cfg := config.LoadConfig()

	r := app.Group("/v1")

	// Public routes
	r.Post("/login", userhandler.Login)
	r.Post("/user", userhandler.CreateUser)
	r.Use(middleware.AuthMiddleware(cfg))

	d := app.Group("/dept")

	d.Post("/dept-create", depthandler.InputDepartment)
	d.Get("/dept-get", depthandler.GetAllDepartment)
	d.Delete("/dept-del", depthandler.DeleteDepartment)
	d.Put("/dept-put/:id", depthandler.UpdateDepartment)



	log.Fatal(app.Listen(":3000"))

}