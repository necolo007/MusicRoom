package router

import (
	"github.com/gin-gonic/gin"
	"github.com/necolo007/MusicRoom/internal/app/user/user_handler"
)

func GenerateRouters(r *gin.Engine) *gin.Engine {
	m1 := r.Group("/api/user")
	{
		m1.POST("/login", user_handler.LoginHandler)
		m1.POST("/register", user_handler.RegisterHandler)
	}
	return r
}
