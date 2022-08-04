package controllers

import (
	"go-fruit-cart/pkg/apperrors"
	"go-fruit-cart/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddProduct(c *gin.Context) {
	email := c.MustGet("claims")

	val, err := models.CheckIfAdmin(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if val {
		theProduct := &models.Product{}
		if err := c.BindJSON(&theProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidInput.Error()})
			c.Abort()
			return
		}
		err = models.InsertNewProduct(theProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "product added"})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
		c.Abort()
		return
	}
}

func GetAllProducts(c *gin.Context) {
	val, err := models.FindAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, val)
}

func DeleteProduct(c *gin.Context) {
	email := c.MustGet("claims")

	val, err := models.CheckIfAdmin(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if val {
		vars := c.Param("id")
		code, err := models.DeleteProductById(vars)
		if err != nil {
			if code == 500 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			} else if code == 404 {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				c.Abort()
				return
			}

		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "successfully deleted"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
		c.Abort()
		return
	}
}

func UpdateProduct(c *gin.Context) {
	email := c.MustGet("claims")

	val, err := models.CheckIfAdmin(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if val {
		vars := c.Param("id")
		theProduct := &models.Product{}
		if err := c.BindJSON(&theProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidInput.Error()})
			c.Abort()
			return
		}
		code, err := models.UpdateProductById(theProduct, vars)
		if err != nil {
			if code == 500 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			} else if code == 404 {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				c.Abort()
				return
			}

		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "product updated successfully"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
		c.Abort()
		return
	}
}
