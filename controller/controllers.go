package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/abdukarimxalilov/ecommerce-go/database"
	"github.com/abdukarimxalilov/ecommerce-go/model"
	generate "github.com/abdukarimxalilov/ecommerce-go/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.UserData(database.Client, "Users")
var prodCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string)(bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true 
	msg := ""

	if err != nil{
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}

func Signup() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)

		defer cancel()

		var user model.User
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return  
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return 
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil{
			log.Panic()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0{
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancel()

		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0{
			c.JSON(http.StatusBadRequest, gin.H{"error":"this phone no. is already in use"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, *&user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]model.ProductUser, 0)
		user.Address_Details = make([]model.Address, 0)
		user.Order_Status = make([]model.Order, 0)

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"the user did not get created"})
			return
		}

		defer cancel()

		c.JSON(http.StatusCreated, "Successfully signed in!")
	}
}

func Login() gin.HandlerFunc{
	return func (c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		var user model.User
		var founduser model.User
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return 
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)

		defer cancel()

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)

		c.JSON(http.StatusFound, founduser)
	}
}


func ProductViewerAdmin() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products model.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.Product_ID = primitive.NewObjectID()
		_, anyerr := prodCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
	}
}


func SearchProduct() gin.HandlerFunc{
	return func(c *gin.Context) {

		var productlist []model.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		cursor, err := prodCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong, please try after some time")
			return
		}

		cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return 
		}
	
	defer cursor.Close(ctx)

	if err := cursor.Err(); err != nil{
		log.Println(err)
		c.IndentedJSON(400, "invalid")
	}
	defer cancel()
	c.IndentedJSON(200, productlist)
	}
}

func SearchProductByQuery() gin.HandlerFunc{
	return func(c *gin.Context) {
		var searchProducts []model.Product
		queryParam := c.Query("name")

		if queryParam == ""{
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		searchqueryDB, err := prodCollection.Find(ctx, bson.M{"product name": bson.M{"$regex":queryParam}}) 

		if err != nil{
			c.IndentedJSON(404, "something went wrong while fetching the data")
			return
		}

		searchqueryDB.All(ctx, &searchProducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer searchqueryDB.Close(ctx)


		if err := searchqueryDB.Err(); err != nil{
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return 
		}

		defer cancel()
		c.IndentedJSON(200, searchProducts)
		
	}
}