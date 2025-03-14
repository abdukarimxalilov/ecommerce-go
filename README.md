# Ecommerce-Go

**E-Commerce API with Golang, Gin, and MongoDB**

## Overview

Ecommerce-Go is a backend API for an e-commerce platform, built using Golang with the Gin framework and MongoDB. It is containerized using Docker to ensure a consistent development environment.

## Features

- **User Authentication** (Signup, Login)
- **Product Management** (Add, View, Search)
- **Cart Management** (Add, Remove, Checkout, Instant Buy)
- **Address Management** (Add, Edit, Delete)
- **JWT Authentication** for security
- **Docker Compose** for MongoDB setup

## Project Structure

```
├── controller      # Business logic for users, products, cart, etc.
├── database        # MongoDB setup and collections
├── middleware      # JWT authentication logic
├── model           # Data models for MongoDB
├── routes          # API route definitions
├── tokens          # JWT token generation and validation
├── main.go         # Entry point of the application
├── docker-compose.yml # Docker configuration for MongoDB
```

## Setup and Installation

### Clone the Repository

```sh
git clone https://github.com/abdukarimxalilov/ecommerce-go.git
cd ecommerce-go
```

### Initialize the Go Module

```sh
go mod init github.com/abdukarimxalilov/ecommerce-go
go mod tidy  # Install dependencies
```

### Start MongoDB Using Docker

```sh
docker-compose up -d  # Runs MongoDB in a container
```

> **Note:** No need to install MongoDB manually; Docker handles it.

### Run the API Server

```sh
go run main.go
```

## API Endpoints

### User Routes

```go
func UserRoutes(router *gin.Engine) {
    router.POST("/users/signup", controller.Signup)
    router.POST("/users/login", controller.Login)
}
```

### Product Routes

```go
func ProductRoutes(router *gin.Engine) {
    router.POST("/admin/addproduct", controller.ProductViewerAdmin)
    router.GET("/users/productview", controller.SearchProduct)
    router.GET("/users/search", controller.SearchProductByQuery)
}
```

### Cart Routes

```go
func CartRoutes(router *gin.Engine) {
    router.GET("/addtocart", app.AddToCart)
    router.GET("/removeitem", app.RemoveItem)
    router.GET("/cartcheckout", app.BuyFromCart)
    router.GET("/instantbuy", app.InstantBuy)
    router.GET("/listcart", controller.GetItemFromCart)
}
```

### Address Management

```go
func AddressRoutes(router *gin.Engine) {
    router.POST("/addaddress", controller.AddAddress)
    router.PUT("/edithomeaddress", controller.EditHomeAddress)
    router.PUT("/editworkaddress", controller.EditWorkAddress)
    router.GET("/deleteaddresses", controller.DeleteAddress)
}
```

## Docker Compose Configuration (`docker-compose.yml`)

```yaml
version: '3.1'
services:
  mongo:
    image: mongo:5.0.3
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: development
      MONGO_INITDB_ROOT_PASSWORD: testpassword
  
  mongo-express:
    image: mongo-express
    ports:
      - "8081:8081"
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: development
      ME_CONFIG_MONGODB_ADMINPASSWORD: testpassword
      ME_CONFIG_MONGODB_URL: mongodb://development:testpassword@mongo:27017/
```

> **Note:** MongoDB is containerized, so no manual installation is required.

## Development Status

This project is under active development. More features and improvements will be added over time. Contributions and feedback are welcome!

### Stop Docker Containers

```sh
docker-compose down
```

## 🤝 Contributing

1. Fork the repository.
2. Create a new branch:
   ```sh
   git checkout -b feature-name
   ```
3. Commit changes:
   ```sh
   git commit -m 'Add new feature'
   ```
4. Push to the branch:
   ```sh
   git push origin feature-name
   ```
5. Open a Pull Request.

