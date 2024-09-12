package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

const cartIdParamEndpoint = "api/carts/:id"
const productIdParamEndpoint = "api/products/:id"
const categoriesIdParamEndpoint = "api/categories/:id"

var oauthStateString = "randomState"

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "REDIRECT_URL/callback",
	ClientID:     "GOOGLE_CLIENT_ID",
	ClientSecret: "GOOGLE_CLIENT_SECRET",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

// Session store for handling user sessions
var store = sessions.NewCookieStore([]byte("secret-key"))

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://88d1-2a02-a31b-20c9-cd80-c534-5962-9761-2df4.ngrok-free.app"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))
	initDB()

	e.POST("api/products", CreateProduct)
	e.GET("api/products", getAllProducts)
	e.GET(productIdParamEndpoint, getProduct)
	e.PUT(productIdParamEndpoint, updateProduct)
	e.DELETE(productIdParamEndpoint, deleteProduct)
	e.POST("api/carts", createCart)
	e.GET(cartIdParamEndpoint, getCart)
	e.POST("api/carts/add", addToCart)
	e.PUT(cartIdParamEndpoint, updateCart)
	e.DELETE(cartIdParamEndpoint, deleteCart)
	e.POST("api/categories", createCategory)
	e.GET(categoriesIdParamEndpoint, getCategory)
	e.PUT(categoriesIdParamEndpoint, updateCategory)
	e.DELETE(categoriesIdParamEndpoint, deleteCategory)
	e.GET("api/categories/:category_id/products", getProductsByCategory)
	e.POST("api/payments", createPayment)
	e.GET("api/payments/:id", getPayment)

	e.GET("api/login", handleGoogleLogin)
	e.GET("api/callback", handleGoogleCallback)
	e.GET("api/profile", handleProfile)
	e.GET("api/logout", handleLogout)

	e.Logger.Fatal(e.Start(":8080"))
}

func initDB() {
	// Use environment variables for the database connection
	host := "localhost" //os.Getenv("POSTGRES_HOST")
	user := "postgres"  //os.Getenv("POSTGRES_USER")
	password := "DB_PASS"
	dbname := "DB_NAME"
	port := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	// Assign the DB connection to the global `db` variable
	var err error // Declare err separately to avoid shadowing the global `db` variable
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initialize database, got error: %v", err)
	}
	fmt.Println("Connected to the database:", db)
}
