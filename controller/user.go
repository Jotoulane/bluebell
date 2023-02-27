package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	// 1.获取参数，参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"mag": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errors.Translate(trans)), //翻译为中文的错误
		})
		return
	}
	fmt.Printf("p%v\n", p)
	// 2.业务处理
	if err := logic.SignUp(p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}
	// 3.返回相应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func LoginHandler(c *gin.Context) {
	//获取请求参数，参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"mag": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errors.Translate(trans)), //翻译为中文的错误
		})
		return
	}
	//业务逻辑处理
	if err := logic.Login(p); err != nil {
		zap.L().Error("logic.Login failed", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户名或者密码错误",
		})
		return
	}
	//返回相应
	c.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
	})
}
