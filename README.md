# ecommerce-go

E-Commerce API with Golang, Gin, and MongoDB

# Overview

This is an e-commerce backend API built using Golang with the Gin framework and MongoDB as the database. The project is containerized using Docker to simplify database setup and ensure a consistent development environment.

# Features

User Authentication (Signup, Login)

Product Management (Add, View, Search)

Cart Management (Add, Remove, Checkout, Instant Buy)

Address Management (Add, Edit, Delete)

JWT Authentication for security

Docker Compose for MongoDB setup

# Project Structure

- ├── controller     # Business logic (APIs for user, product, cart, etc.)
- ├── database       # MongoDB setup and collections
- ├── middleware     # JWT authentication logic
- ├── model          # Data models for MongoDB
- ├── routes         # API route definitions
- ├── tokens         # JWT token generation and validation
- ├── main.go        # Entry point of the application
- ├── docker-compose.yml  # Docker configuration for MongoDB

# Setup and Installation

1. Clone the Repository

git clone https://github.com/your-username/your-repo.git
cd your-repo

2. Initialize the Go Module

go mod init github.com/your-username/your-repo
go mod tidy  # Install dependencies

3. Start MongoDB Using Docker

docker-compose up -d  # Runs MongoDB in a container

Note: You don't need to install MongoDB manually; Docker will handle it.

4. Run the API Server

go run main.go

API Endpoints

# User Routes

func UserRoutes(incomingRoutes *gin.Engine) {
    incomingRoutes.POST("/users/signup", controller.Signup)
    incomingRoutes.POST("/users/login", controller.Login)
}

# Product Routes

func ProductRoutes(incomingRoutes *gin.Engine) {
    incomingRoutes.POST("/admin/addproduct", controller.ProductViewerAdmin)
    incomingRoutes.GET("/users/productview", controller.SearchProduct)
    incomingRoutes.GET("/users/search", controller.SearchProductByQuery)
}

# Cart Routes

func CartRoutes(router *gin.Engine) {
    router.GET("/addtocart", app.AddToCart)
    router.GET("/removeitem", app.RemoveItem)
    router.GET("/cartcheckout", app.BuyFromCart)
    router.GET("/instantbuy", app.InstantBuy)
    router.GET("/listcart", controller.GetItemFromCart)
}

# Address Management

func AddressRoutes(router *gin.Engine) {
    router.POST("/addaddress", controller.AddAddress)
    router.PUT("/edithomeaddress", controller.EditHomeAddress)
    router.PUT("/editworkaddress", controller.EditWorkAddress)
    router.GET("/deleteaddresses", controller.DeleteAddress)
}

# Docker Compose Configuration

docker-compose.yml

version: '3.1'
services:
  mongo:
    image: mongo:5.0.3
    ports:
      - 27017:27017
    environment:
      MONGO_INITDB_ROOT_USERNAME: development
      MONGO_INITDB_ROOT_PASSWORD: testpassword
  
  mongo-express:
    image: mongo-express
    ports:
      - 8081:8081
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: development
      ME_CONFIG_MONGODB_ADMINPASSWORD: testpassword
      ME_CONFIG_MONGODB_URL: mongodb://development:testpassword@mongo:27017/

# Notes

- MongoDB is containerized: You don’t need to install it manually.

- JWT authentication is implemented in the middleware folder.

- Controller functions handle business logic, while routes define API endpoints.

- To stop Docker containers, run:

docker-compose down

# Development Status

This project is still in the development stage. More features, improvements, and documentation updates will be added over time. Contributions and feedback are welcome!



