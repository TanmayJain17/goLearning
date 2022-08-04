package controllers

import (
	"fmt"
	"go-fruit-cart/pkg/apperrors"
	"go-fruit-cart/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCart(c *gin.Context) {
	email := c.MustGet("claims")
	val, err := models.CheckIfAdmin(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if !val {
		/* cartid := c.Param("cartid") */
		data, err := models.GetProductsInCart( /* cartid, */ email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		fmt.Println(string(data))
		c.Data(200, "application/json", data)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
		c.Abort()
		return
	}
}

func AddToCart(c *gin.Context) {
	email := c.MustGet("claims")
	val, err := models.CheckIfAdmin(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if !val {
		/* cartid := c.Param("cartid") */
		productid := c.Param("productid")

		//check if the product in the product list
		name, price, err := models.FindProduct(productid)
		fmt.Println("name=>", name, price)
		if price != 0 {
			if err != nil {
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					c.Abort()
					return
				}
			}

			//now f the product is in the product list then update the cart
			_, err := models.AddProductToCart( /* cartid, */ email, name, price, productid)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			}

			c.IndentedJSON(http.StatusOK, gin.H{"message": "product added"})
			/* c.Data(200, "application/json", data) */
			//c.IndentedJSON(http.StatusOK, data)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrProductNotFound.Error()})
			c.Abort()
			return
		}

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
		c.Abort()
		return
	}
}

func RemoveFromCart(c *gin.Context) {
	email := c.MustGet("claims")
	val, err := models.CheckIfAdmin(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if !val {
		/* cartid := c.Param("cartid") */
		productid := c.Param("productid")

		//check if the product in the product list
		name, price, err := models.FindProduct(productid)

		if err != nil {
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
		}
		if price != 0 {

			//now f the product is in the product list then update the cart
			data, err := models.RemoveProductFromCart( /* cartid, */ email, name, price, productid)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			if string(data) == "null" {
				c.IndentedJSON(http.StatusOK, gin.H{"message": "cart empty"})
			} else {
				c.IndentedJSON(http.StatusOK, gin.H{"message": "product removed"})
			}

			/* c.Data(200, "application/json", "product removed") */
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrProductNotFound.Error()})
			c.Abort()
			return
		}

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
		c.Abort()
		return
	}
}
