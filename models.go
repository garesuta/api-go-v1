package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CompanyInfo struct {
	Name               string `json:"name"`
	CompanyDescription string `json:"Company_description"`
	ContactEmail       string `json:"contactEmail"`
	ContactPhone       string `json:"contactPhone"`
}
type Jobs struct {
	ID          uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string      `json:"title"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Location    string      `json:"location"`
	Salary      string      `json:"salary"`
	Company     CompanyInfo `gorm:"embedded"`
}

// // getJobs retrieves all books
func getJobs(db *gorm.DB, c *fiber.Ctx) error {
	var books []Jobs
	db.Find(&books)
	return c.JSON(books)
}

// // getJobs retrieves a book by id
func getJob(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var book Jobs
	db.First(&book, id)
	return c.JSON(book)
}

// // createBook creates a new job
func createJob(db *gorm.DB, c *fiber.Ctx) error {
	job := new(Jobs)
	if err := c.BodyParser(job); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(job); err != nil {
		return c.Status(fiber.StatusCreated).JSON(job)
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create job"})
}

// // updateJob updates a job by id
func updateJob(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	job := new(Jobs)
	db.First(&job, id)
	if err := c.BodyParser(job); err != nil {
		return err
	}
	db.Save(&job)
	return c.JSON(job)
}

// // deleteBook deletes a book by id
func deleteJob(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&Jobs{}, id)
	return c.SendString("Jobs successfully deleted")
}
