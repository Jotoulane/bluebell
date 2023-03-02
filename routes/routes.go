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

	v1 := r.Group("/api/v1")
	//注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	//登录
	v1.POST("/login", controller.LoginHandler)
	//应用JWT认证中间件
	v1.Use(middleware.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)           //获取社区列表
		v1.GET("/community/:id", controller.CommunityDetailHandler) //根据社区id获取社区详情
		v1.POST("/post", controller.CreatePostHandler)              //创建帖子
		v1.GET("/post/:id", controller.GetPostDetailHandler)        //帖子详情
		v1.GET("/posts/", controller.GetPostListHandler)            //帖子列表
		v1.POST("/vote", controller.PostVoteController)             //帖子投票
	}

	//v1.GET("/ping", middleware.JWTAuthMiddleware(), func(context *gin.Context) {
	//	//如果是登录的用户，判断请求头中是否有 有效的JWT
	//	context.String(http.StatusOK, "pong")
	//})

	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
