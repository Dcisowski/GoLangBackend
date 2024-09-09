package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Category struct {
	gorm.Model
	Name     string    `json:"name"`
	Products []Product `gorm:"many2many:category_products;"`
}

func createCategory(c echo.Context) error {
	var category Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db.Create(&category)
	return c.JSON(http.StatusCreated, category)
}

func getCategory(c echo.Context) error {
	id := c.Param("id")
	var category Category
	if err := db.Preload("Products").First(&category, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Category not found")
	}
	return c.JSON(http.StatusOK, category)
}

func productsByCategory(categoryID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins("JOIN category_products ON category_products.product_id = products.id").
			Where("category_products.category_id = ?", categoryID)
	}
}

func updateCategory(c echo.Context) error {
	id := c.Param("id")
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Cart not found")
	}

	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db.Save(&category)
	return c.JSON(http.StatusOK, category)
}

func deleteCategory(c echo.Context) error {
	id := c.Param("id")
	db.Delete(&Category{}, id)
	return c.NoContent(http.StatusNoContent)
}

func getProductsByCategory(c echo.Context) error {
	categoryIDParam := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDParam, 10, 64) // Parse as unsigned integer
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid category ID")
	}

	var products []Product
	db.Scopes(productsByCategory(uint(categoryID))).Find(&products)
	return c.JSON(http.StatusOK, products)
}
