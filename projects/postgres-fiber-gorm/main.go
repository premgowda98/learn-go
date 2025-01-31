package main

import (
	"log"
	"net/http"
	"os"
	"projects/postgres-fiber-gorm/models"
	"projects/postgres-fiber-gorm/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(ctx *fiber.Ctx) {
	book := Book{}

	err := ctx.BodyParser(&book)

	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "json decode error"})

		return
	}

	err = r.DB.Create(&book).Error

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to create record"})

		return
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "created book",
	})
}

func (r *Repository) GetBooks(ctx *fiber.Ctx) {

	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to get record"})
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "fetched book",
		"data":    bookModels,
	})

}

func (r *Repository) GetBookById(ctx *fiber.Ctx) {

}

func (r *Repository) DeleBook(ctx *fiber.Ctx) {
	bookModel := models.Books{}
	id := ctx.Params("id")

	if id == "" {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "id not found"})
	}

	err := r.DB.Delete(bookModel, id).Error

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete"})
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book deleted",
	})
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/books", r.CreateBook)
	api.Get("/books", r.GetBooks)
	// api.Get("/books/:id", r.GetBookById)
	api.Delete("/books/:id", r.DeleBook)
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal(err)
	}

	err = models.MigrateBooks(db)

	if err != nil {
		log.Fatal(err)
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()

	r.SetupRoutes(app)

	app.Listen(":3000")
}
