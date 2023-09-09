package main

import (
	"context"
	// "fmt"
	"net/http"
	"orderApp/configs"

	// "strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Library struct {
	Book string `json:"book"`
}

var libraryBooks []Library

func main() {
	r := gin.Default()

	config.Connect()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the order App",
		})
	})

	r.POST("/books", func(c *gin.Context) {
		var library Library
		if err := c.ShouldBindJSON(&library); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		collection := config.Client.Database("dOrderApp").Collection("Books")

		result, err := collection.InsertOne(context.Background(), library)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Book added successfully",
			"id":      result.InsertedID,
			"data":    library,
		})
	})

	r.GET("/books", func(c *gin.Context) {
		collection := config.Client.Database("dOrderApp").Collection("Books")
		cursor, err := collection.Find(context.TODO(), bson.D{})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(context.TODO())
		var books []map[string]interface{}

		if err := cursor.All(context.TODO(), &books); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Äll Books Fetched",
			"data":    books,
		})
	})

	r.GET("books/:id", func(c *gin.Context) {
		bookId := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(bookId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ObjectID"})
			return
		}

		collection := config.Client.Database("dOrderApp").Collection("Books")
		filter := bson.M{"_id": objectID}

		var library Library

		if err := collection.FindOne(context.TODO(), filter).Decode(&library); err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Book found",
			"data":    library,
		})
	})

	r.PUT("books/:id", func(c *gin.Context) {
		bookId := c.Param("id")

		objectID, err := primitive.ObjectIDFromHex(bookId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ObjectID"})
			return
		}

		// Create a struct to hold the updated data
		var updatedBook Library

		// Bind the JSON data from the request body to the updatedBook struct
		if err := c.BindJSON(&updatedBook); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		collection := config.Client.Database("dOrderApp").Collection("Books")
		filter := bson.M{"_id": objectID}
		update := bson.M{"$set": updatedBook} // Use updatedBook with the new data

		result, err := collection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Check if the document was found and updated
		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Book updated",
			"data":    updatedBook, // Send the updatedBook data in the response
		})

	})

	r.DELETE("books/:id", func(c *gin.Context) {
		bookId := c.Param("id")

		objectID, err := primitive.ObjectIDFromHex(bookId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ObjectID"})
			return
		}

		collection := config.Client.Database("dOrderApp").Collection("Books")
		filter := bson.M{"_id": objectID}

		result := collection.FindOneAndDelete(context.TODO(), filter)

		c.JSON(http.StatusOK, gin.H{
			"message": "Book deleted",
			"data": result,
		})
	})
	r.Run()
}
