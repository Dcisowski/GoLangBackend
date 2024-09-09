package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

const cartNotFound = "Cart not found"

type Cart struct {
	gorm.Model
	Products []Product `gorm:"many2many:cart_products;"` // Many-to-many relationship between Cart and Product
}

func createCart(c echo.Context) error {
	var cart Cart
	if err := c.Bind(&cart); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db.Create(&cart)
	return c.JSON(http.StatusCreated, cart)
}

type AddToCartRequest struct {
	ProductID int    `json:"product_id"`
	CartID    string `json:"cart_id"`
}

func addToCart(c echo.Context) error {
	req := new(AddToCartRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	productID := req.ProductID // Now we have the productID from the request body
	cartID := req.CartID
	// Optionally handle the cartID
	var product Product

	// Try to find the product by ID
	if err := db.First(&product, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Product not found")
	}
	// Retrieve cart ID from client (e.g., passed via request header or localStorage on frontend)
	var cart Cart

	if cartID == "" {
		// If no cart ID is passed, create a new cart
		cart = Cart{}
		db.Create(&cart)
		cartID = fmt.Sprintf("%d", cart.ID) // Convert cart ID to string for frontend use
	} else {
		// Try to find the existing cart by its ID
		if err := db.First(&cart, cartID).Error; err != nil {
			return c.JSON(http.StatusNotFound, cartNotFound)
		}
	}

	// Add the product to the existing or new cart
	db.Model(&cart).Association("Products").Append(&product)

	// Return the updated cart, including the cart ID so the frontend can store it
	return c.JSON(http.StatusOK, map[string]interface{}{
		"cart_id": cart.ID,
		"cart":    cart,
	})
}

func getCart(c echo.Context) error {
	cartID := c.Param("id")
	var cart Cart
	// Fetch cart with associated products
	if err := db.Preload("Products").First(&cart, cartID).Error; err != nil {
		return c.JSON(http.StatusNotFound, cartNotFound)
	}

	// Calculate the total price
	total := 0.0
	for _, product := range cart.Products {
		total += product.Price
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"cart":  cart,
		"total": total,
	})
}

func updateCart(c echo.Context) error {
	id := c.Param("id")
	var cart Cart
	if err := db.First(&cart, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, cartNotFound)
	}

	if err := c.Bind(&cart); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db.Save(&cart)
	return c.JSON(http.StatusOK, cart)
}

func deleteCart(c echo.Context) error {
	id := c.Param("id")
	db.Delete(&Cart{}, id)
	return c.NoContent(http.StatusOK)
}
