package main

import (
	"log"
	"os"

	"github.com/abdukarimxalilov/ecommerce-go/controller"
	"github.com/abdukarimxalilov/ecommerce-go/database"
	"github.com/abdukarimxalilov/ecommerce-go/middleware"
	"github.com/abdukarimxalilov/ecommerce-go/routes"

	"github.com/gin-gonic/gin"
)

func main(){
	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	app := controller.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	
	router.GET("/listcart", controller.GetItemFromCart())
	router.POST("/addaddress", controller.AddAddress())
	router.PUT("/edithomeaddress", controller.EditHomeAddress())
	router.PUT("/editworkaddress", controller.EditWorkAddress())
	router.GET("/deleteaddresses", controller.DeleteAddress())
	log.Fatal(router.Run(":" + port))

}