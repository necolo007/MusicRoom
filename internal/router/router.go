package router

import (
	"github.com/gin-gonic/gin"
	"github.com/necolo007/MusicRoom/core/middleware/web"
	"github.com/necolo007/MusicRoom/internal/app/music/music_handler"
	"github.com/necolo007/MusicRoom/internal/app/user/user_handler"
)

func GenerateRouters(r *gin.Engine) *gin.Engine {
	m1 := r.Group("/api/user")
	{
		m1.POST("/login", user_handler.LoginHandler)
		m1.POST("/register", user_handler.RegisterHandler)
	}
	m2 := r.Group("/api/music", web.JWTAuthMiddleware())
	{
		m2.POST("/upload", music_handler.UploadMusic)
		m2.GET("/list", music_handler.ListMusic)
		m2.GET("/music-detail", music_handler.GetMusic)
		m2.DELETE("/delete", music_handler.DeleteMusic)
	}
	return r
}
