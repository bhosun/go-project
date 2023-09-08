package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Library struct {
	ID   int    `json:"id"`
	Book string `json:"book"`
}

var libraryBooks []Library

func main() {
	r := gin.Default()

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

		library.ID = len(libraryBooks) + 1
		libraryBooks = append(libraryBooks, library)
		c.JSON(http.StatusCreated, libraryBooks)
	})

	r.GET("/books", func(c *gin.Context) {
		c.JSON(http.StatusOK, libraryBooks)
	})

	r.GET("books/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, book := range libraryBooks {
			if id == strconv.Itoa(book.ID) {
				c.JSON(http.StatusOK, book)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
	})

	r.Run()
}
