package controller

import (
	"bluebell/logic"
	models "bluebell/models"
	"strconv"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	//获取参数以及校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("CreatePost with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errors.Translate(trans)))
		return
	}
	//从c取到当前发送请求的用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//创建帖子
	if err = logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 根据帖子的id获取帖子的详情
func GetPostDetailHandler(c *gin.Context) {
	//获取参数
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//根据id取出帖子数据
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回相应
	ResponseSuccess(c, data)
}

// GetPostListHandler 帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {
	pageNum, pageSize := getPostList(c)
	data, err := logic.GetPostList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// PostVoteController 为帖子投票
func PostVoteController(c *gin.Context) {
	//参数校验
	vote := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(vote); err != nil {
		//请求参数有问题
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	userId, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//投票业务逻辑
	if err = logic.VoteForPost(userId, vote); err != nil {
		zap.L().Error("logic.VoteForPost")
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// GetPostListHandler2 接口帖子列表升级版
// 根据前端的数据按照时间或者分数排序
func GetPostListHandler2(c *gin.Context) {
	//获取请求参数
	p := models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//去redis查询id列表
	data, err := logic.GetPostListNew(&p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

//// GetCommunityPostListHandler 根据社区查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	//获取请求参数
//	p := &models.ParamCommunityPostList{
//		ParamPostList: &models.ParamPostList{Size: 10, Order: models.OrderTime},
//	}
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//	//去redis查询id列表
//
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("logic.GetCommunityPostList() failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	ResponseSuccess(c, data)
//}
