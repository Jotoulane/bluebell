package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUp(mode string) *gin.Engine {

	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成为发布模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//注册业务路由
	r.POST("/signup", controller.SignUpHandler)
	//登录
	r.POST("/login", controller.LoginHandler)

	r.GET("/ping", middleware.JWTAuthMiddleware(), func(context *gin.Context) {
		//如果是登录的用户，判断请求头中是否有 有效的JWT
		context.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
