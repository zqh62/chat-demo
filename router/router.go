package router

import "github.com/gin-gonic/gin"
import "chat/api"

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery(), gin.Logger())
	//Recovery会恢复所有panic
	//Logger日志
	v1 := r.Group("/")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
		v1.POST("user/register", api.UserRegister)
	}
	return r
}
