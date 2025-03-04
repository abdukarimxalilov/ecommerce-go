package routes

import(
	"github.com/abdukarimxalilov/ecommerce-go/controller"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/users/signup", controller.Signup())
	incomingRoutes.POST("/users/login", controller.Login())
	incomingRoutes.POST("/admin/addproduct", controller.ProductViewerAdmin())
	incomingRoutes.GET("/users/productview", controller.SearchProduct())
	incomingRoutes.GET("/users/search", controller.SearchProductByQuery())
}