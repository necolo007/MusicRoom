package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/necolo007/MusicRoom/core/auth"
	"github.com/necolo007/MusicRoom/core/gin/dbs"
	"github.com/necolo007/MusicRoom/internal/router"
)

func GinInit() *gin.Engine {
	r := gin.Default()
	dbs.InitDB()
	router.GenerateRouters(r)

	auth.InitSecret()
	return r
}
