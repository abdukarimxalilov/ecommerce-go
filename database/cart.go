package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/abdukarimxalilov/ecommerce-go/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

var(
	ErrCantFindProduct = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't find the product")
	ErrUserIdIsNotValid = errors.New("this user is not valid")
	ErrCantUpdateUser = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem = errors.New("was unable to get item from the cart")
	ErrCantBuyCartItem = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error{
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var productcart []model.ProductUser
	err = searchfromdb.All(ctx, &productcart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productcart}}}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartFrom(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error{
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil{
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull":bson.M{"usercart": bson.M{"_id":productID}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil{
		return ErrCantRemoveItemCart
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error{

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil{
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getcartitems model.User
	var ordercart model.Order

	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Order_At = time.Now()
	ordercart.Order_Cart = make([]model.ProductUser, 0)
	ordercart.Payment_Method.COD = true 

	unwind := bson.D{{Key:"$unwind", Value:bson.D{primitive.E{Key:"path", Value: "$usercart"}}}}
	grouping := bson.D{{Key:"$group", Value:bson.D{primitive.E{Key:"_id", Value:"$_id"}, {Key:"total", Value:bson.D{primitive.E{Key:"$sum", Value:"$usercart.price"}}}}}}
	currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil{
		panic(err)
	}

	var getusercart []bson.M 
	if err = currentresults.All(ctx, &getusercart); err != nil{
		panic(err)
	}

	var total_price int32 
	for _, user_item := range getusercart{
		price := user_item["total"]
		total_price = price.(int32)
	}
	ordercart.Price = int(total_price)

	filter := bson.D{primitive.E{Key:"_id", Value: id}}
	update := bson.D{{Key:"$push", Value:bson.D{primitive.E{Key:"orders", Value:ordercart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil{
		log.Println(err)
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key:"_id", Value:id}}).Decode(&getcartitems)
	if err != nil{
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push":bson.M{"orders.$[].order_list":bson.M{"$each":getcartitems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil{
		log.Println(err)
	}

	usercart_empty := make([]model.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key:"$set", Value:bson.D{primitive.E{Key:"usercart", Value:usercart_empty}}}}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil{
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error{
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil{
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var product_details model.ProductUser
	var order_detail model.Order 

	order_detail.Order_ID = primitive.NewObjectID()
	order_detail.Order_At = time.Now()
	order_detail.Order_Cart = make([]model.ProductUser, 0)
	order_detail.Payment_Method.COD = true 
	prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details) 
	if err != nil{
		log.Println(err)
	}

	order_detail.Price = product_details.Price 

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key:"$push", Value:bson.D{primitive.E{Key:"order_details", Value:order_detail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list":product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}

