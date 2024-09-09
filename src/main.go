package main

import (
	"fmt"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

const cartIdParamEndpoint = "/carts/:id"
const productIdParamEndpoint = "/products/:id"
const categoriesIdParamEndpoint = "/categories/:id"

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))
	initDB()

	e.POST("/products", CreateProduct)
	e.GET("/products", getAllProducts)
	e.GET(productIdParamEndpoint, getProduct)
	e.PUT(productIdParamEndpoint, updateProduct)
	e.DELETE(productIdParamEndpoint, deleteProduct)
	e.POST("/carts", createCart)
	e.GET(cartIdParamEndpoint, getCart)
	e.POST("/carts/add", addToCart)
	e.PUT(cartIdParamEndpoint, updateCart)
	e.DELETE(cartIdParamEndpoint, deleteCart)
	e.POST("/categories", createCategory)
	e.GET(categoriesIdParamEndpoint, getCategory)
	e.PUT(categoriesIdParamEndpoint, updateCategory)
	e.DELETE(categoriesIdParamEndpoint, deleteCategory)
	e.GET("/categories/:category_id/products", getProductsByCategory)
	e.POST("/payments", createPayment)
	e.GET("/payments/:id", getPayment)

	e.Logger.Fatal(e.Start(":8080"))
}

func initDB() *gorm.DB {
	// Use environment variables for the database connection
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initialize database, got error: %v", err)
	}
	fmt.Println("Connected to the database:", db)
	return db
}
