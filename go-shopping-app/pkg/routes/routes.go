package routes

import (
	"go-fruit-cart/pkg/controllers"
	"go-fruit-cart/pkg/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func HandleAllRoutes(router *gin.Engine) {
	router.POST("/login", controllers.UserLogin)
	router.POST("/signup", controllers.NewUser)
	router.GET("/products", controllers.GetAllProducts)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	secured := router.Group("/auth").Use(middleware.Auth())
	{
		/* secured.GET("/ping", controllers.Ping) */
		secured.POST("/product", controllers.AddProduct)
		secured.DELETE("/product/:id", controllers.DeleteProduct)
		secured.PUT("/product/:id", controllers.UpdateProduct)

		secured.GET("/cart", controllers.GetCart)

		secured.PUT("/products/:productid/cart", controllers.AddToCart)
		/* secured.PUT("/carts/:cartid/products/:productid/remove", controllers.RemoveFromCart) */
		secured.DELETE("/products/:productid/cart", controllers.RemoveFromCart)

	}

}
