package main

import (
	"log"

	"hotelmanagementsystem.com/v2/config"
	"hotelmanagementsystem.com/v2/database"
	"hotelmanagementsystem.com/v2/handlers"
	"hotelmanagementsystem.com/v2/middleware"
	"hotelmanagementsystem.com/v2/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DBDriver, cfg.DatabaseURL, cfg.SQLitePath)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Guest{},
		&models.Room{},
		&models.Booking{},
		&models.HotelSettings{},
		&models.IncomeCategory{},
		&models.IncomeRecord{},
		&models.ExpenditureRecord{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := database.Seed(db, cfg); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "Hotel Management System API",
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())
	app.Static("/uploads", "./uploads")

	api := app.Group("/api")
	api.Get("/health", handlers.Health)

	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret, cfg.TokenTTLHours)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	protected := api.Group("", middleware.Protected(cfg.JWTSecret))

	dashboardHandler := handlers.NewDashboardHandler(db)
	roomHandler := handlers.NewRoomHandler(db)
	guestHandler := handlers.NewGuestHandler(db)
	bookingHandler := handlers.NewBookingHandler(db)
	settingsHandler := handlers.NewSettingsHandler(db)
	incomeCategoryHandler := handlers.NewIncomeCategoryHandler(db)
	incomeHandler := handlers.NewIncomeHandler(db)
	expenditureHandler := handlers.NewExpenditureHandler(db)

	protected.Get("/dashboard", dashboardHandler.GetSummary)

	protected.Get("/rooms", roomHandler.ListRooms)
	protected.Post("/rooms", roomHandler.CreateRoom)
	protected.Patch("/rooms/:id/status", roomHandler.UpdateRoomStatus)

	protected.Get("/guests", guestHandler.ListGuests)
	protected.Post("/guests", guestHandler.CreateGuest)

	protected.Get("/bookings", bookingHandler.ListBookings)
	protected.Post("/bookings", bookingHandler.CreateBooking)
	protected.Patch("/bookings/:id/status", bookingHandler.UpdateBookingStatus)

	protected.Get("/settings/hotel", settingsHandler.GetHotelSettings)
	protected.Put("/settings/hotel", settingsHandler.UpdateHotelSettings)

	protected.Get("/income-categories", incomeCategoryHandler.ListCategories)
	protected.Post("/income-categories", incomeCategoryHandler.CreateCategory)

	protected.Get("/incomes", incomeHandler.ListIncomeRecords)
	protected.Post("/incomes", incomeHandler.CreateIncomeRecord)

	protected.Get("/expenditures", expenditureHandler.ListExpenditureRecords)
	protected.Post("/expenditures", expenditureHandler.CreateExpenditureRecord)

	log.Printf("hotel api running on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
