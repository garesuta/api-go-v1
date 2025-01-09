package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// "github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_HOST")), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect to database")
	}
	// Migrate the schema
	db.AutoMigrate(&CompanyInfo{}, &Jobs{})

	// Setup Fiber
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// // CRUD routes
	app.Get("/jobs", func(c *fiber.Ctx) error {
		return getJobs(db, c)
	})
	app.Get("/jobs/:id", func(c *fiber.Ctx) error {
		return getJob(db, c)
	})
	app.Post("/jobs", func(c *fiber.Ctx) error {
		return createJob(db, c)
	})
	app.Put("/jobs/:id", func(c *fiber.Ctx) error {
		return updateJob(db, c)
	})
	app.Delete("/jobs/:id", func(c *fiber.Ctx) error {
		return deleteJob(db, c)
	})

	// Start server
	log.Fatal(app.Listen(":8000"))
}
