package order

import (
	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine) {
	r.POST("/books", AddOrder())
	r.GET("/books", GetBooks())
	r.GET("/books/:id", GetBookById())
	r.PUT("books/:id", EditBook())
	r.DELETE("books/:id", DeleteBook())
}