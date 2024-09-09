package main

import (
	_ "fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Initialize Echo instance
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))
	// Initialize database connection
	initDB()

	// Routes for Product CRUD
	e.POST("/products", createProduct)
	e.GET("/products", getAllProducts)
	e.GET("/products/:id", getProduct)
	e.PUT("/products/:id", updateProduct)
	e.DELETE("/products/:id", deleteProduct)
	e.POST("/carts", createCart)
	e.GET("/carts/:id", getCart)
	e.POST("/carts/add", addToCart)
	e.PUT("/carts/:id", updateCart)
	e.DELETE("/carts/:id", deleteCart)
	e.POST("/categories", createCategory)
	e.GET("/categories/:id", getCategory)
	e.PUT("/categories/:id", updateCategory)
	e.DELETE("/categories/:id", deleteCategory)
	e.GET("/categories/:category_id/products", getProductsByCategory)

	e.Logger.Fatal(e.Start(":8080"))
}

func initDB() {
	var err error

	// Define the DSN (Data Source Name) for PostgreSQL
	dsn := "host=localhost user=user password=password dbname=dbname port=5432 sslmode=disable TimeZone=Europe/Warsaw"

	// Open the connection using GORM and the PostgreSQL driver
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Auto-migrate your models
	err = db.AutoMigrate(&Product{}, &Cart{}, &Category{})
	if err != nil {
		return
	}
}
