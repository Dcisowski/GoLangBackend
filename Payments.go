package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

// Payment model
type Payment struct {
	gorm.Model
	CartID        uint    `json:"cart_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	Status        string  `json:"status"` // Example values: "pending", "completed", "failed"
}

// Endpoint to handle payment creation
func createPayment(c echo.Context) error {
	var payment Payment
	if err := c.Bind(&payment); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db.Create(&payment)
	return c.JSON(http.StatusCreated, payment)
}

// Endpoint to get payment by ID
func getPayment(c echo.Context) error {
	id := c.Param("id")
	var payment Payment
	if err := db.First(&payment, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Payment not found")
	}
	return c.JSON(http.StatusOK, payment)
}
