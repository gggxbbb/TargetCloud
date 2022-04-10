package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	err := r.Run()
	if err != nil {
		return
	}
}
