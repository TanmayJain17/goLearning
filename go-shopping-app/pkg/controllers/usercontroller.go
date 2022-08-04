package controllers

import (
	"fmt"
	"go-fruit-cart/pkg/apperrors"
	"go-fruit-cart/pkg/models"
	"go-fruit-cart/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewUser(c *gin.Context) {
	theUser := &models.User{}
	if err := c.BindJSON(&theUser); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": apperrors.ErrInvalidInput.Error()})
		c.Abort()
		return
	}
	/* err := utils.ValidateEmail(theUser.Email) */
	/* if err == nil { */
	hashedPassword, _ := utils.HashPassword(theUser.Password)
	/* if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	} */
	theUser.Password = hashedPassword
	code, err := models.InsertNewUser(theUser)
	if err != nil {
		if code == 502 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		} else if code == 401 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "SignUp Successfull"})

	/* } else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	} */
}

func UserLogin(c *gin.Context) {
	var request models.TokenRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidInput.Error()})
		c.Abort()
		return
	}
	firstname, lastname, _, code, err := models.FindTheUser(&request)
	if err == nil {

		tokenString, _ := utils.GenerateJWTToken(request.Email, firstname, lastname)
		/* if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		} */
		fmt.Println(tokenString)
		/* if isadmin {
			adminToken = tokenString
			adminTokenAd = &adminToken
		} else {
			userToken = tokenString
			userTokenAd = &userToken
		} */

		c.JSON(http.StatusOK, gin.H{"message": "successfully logged in"})
	} else {
		if code == 502 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		} else if code == 404 {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrInvalidCredentials.Error()})
			c.Abort()
			return
		}

	}
}
