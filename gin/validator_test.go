package main

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestValidator(t *testing.T) {
	server := gin.Default()
	server.POST("/users/signup", func(ctx *gin.Context) {
		var u User
		if err := ctx.Bind(&u); err != nil {
			t.Log(err)
		}
	})
	server.Run(":8083")
}

type User struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}
