package router

import (
	"log"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/config"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/handler"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/middleware"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App, userhandler *handler.UserHandler, depthandler *handler.DeptHandler, authHandler *handler.AuthHandler, lokasiHandler *handler.LokasiHandler) {
	cfg := config.LoadConfig()

	a := app.Group("/auth")

	a.Post("/login", authHandler.Login)
	a.Post("/change-password", middleware.AuthMiddleware(cfg), authHandler.ChangePassword)
	a.Post("/logout", middleware.AuthMiddleware(cfg), authHandler.Logout)

	r := app.Group("/v1")

	r.Post("/user", userhandler.CreateUser)
	r.Get("/user/search", userhandler.GetUsers)
	r.Put("/user-put/:id", userhandler.UpdateUserHandler)
	r.Delete("/user-del", userhandler.DeleteUserHandler)
	r.Use(middleware.AuthMiddleware(cfg))

	d := app.Group("/dept")

	d.Post("/dept-create", depthandler.InputDepartment)
	d.Get("/dept-get", depthandler.GetAllDepartment)
	d.Delete("/dept-del", depthandler.DeleteDepartment)
	d.Put("/dept-put/:id", depthandler.UpdateDepartment)

	l := app.Group("/lok")

	l.Post("/lokasi", lokasiHandler.CreateLokasi)
	l.Get("/lokasi-get", lokasiHandler.GetAllLocationsWithDetails)
	l.Put("/lokasi-put/:id", lokasiHandler.UpdateLokasi)
	l.Delete("/lokasi-del", lokasiHandler.DeleteLokasi)
	l.Put("/detail-put/:id", lokasiHandler.UpdateDetailLokasi)
	l.Delete("/detail-del", lokasiHandler.DeleteDetailLokasi)


	log.Fatal(app.Listen(":3000"))

}