package router

import (
	"log"

	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/config"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/handler"
	"github.com/SKARIGA-RPL-XII/project-kik-farfan02/Roomify/roomify/middleware"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App, 
	userhandler *handler.UserHandler, 
	depthandler *handler.DeptHandler, 
	authHandler *handler.AuthHandler, 
	lokasiHandler *handler.LokasiHandler, 
	settingHandler *handler.SettingHandler,
	bookinghandler *handler.BookingHandler) {
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

	s := app.Group("/set")

	s.Get("/all", settingHandler.GetAllSettings)
	s.Get("/get-setting", settingHandler.GetSetting)
	s.Put("/put-setting", settingHandler.UpdateSetting)

	b := app.Group("/book")

	b.Post("/booking-create", middleware.AuthMiddleware(cfg), bookinghandler.CreateBooking)
	b.Get("/getbyid", bookinghandler.GetBookingByID)
	b.Get("/getby-user",middleware.AuthMiddleware(cfg), bookinghandler.GetBookingsByUser)
	b.Get("/getall", bookinghandler.GetAllBookings)
	b.Put("/booking-put/:id", bookinghandler.UpdateBooking)
	b.Delete("/booking-del",middleware.AuthMiddleware(cfg), bookinghandler.DeleteBooking)

	log.Fatal(app.Listen(":3000"))

}