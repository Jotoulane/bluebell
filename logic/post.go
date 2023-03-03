package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"strconv"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	//生成postID
	p.ID = snowflake.GenID()
	//保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		zap.L().Error("mysql.CreatePost(p) failed", zap.Int64("p.ID", p.ID), zap.Error(err))
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	if err != nil {
		zap.L().Error("redis.CreatePost(p.ID) failed", zap.Int64("p.ID", p.ID), zap.Error(err))
		return err
	}
	return
}

// GetPostById 根据帖子id查询帖子详情
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	//查询数据，组合接口数据
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
		return
	}

	//根据社区信息查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Int64("CommunityID", post.CommunityID), zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(pageNum int64, PageSize int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(pageNum, PageSize)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
			continue
		}

		//根据社区信息查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Int64("CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2 获取帖子列表
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//redis查询id列表
	ids, err := redis.GetPostIDInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	//mysql查询帖子详情
	//返回的数据还要按照我给定的id数据进行返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//查询每篇帖子的投票信息
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将帖子作者以及分区信息查询出来填充到帖子当中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
			continue
		}

		//根据社区信息查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Int64("CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// VoteForPost 为帖子投票
/*投票的几种情况:
direction=1时，有两种情况:
	1.之前没有投过票，现在投赞成票
	2.之前投反对票,现在改投赞成票
direction=0时,有两种情况:
	1.之前投过赞成票,现在要取消投票
	2.之前投过反对票，现在要取消投票
direction=-1时，有两种情况:
	1.之前没有投过票,现在投反对票
	2.之前投赞成票，现在改投反对票


投票的限制
	每个帖子自发表之日起一个星期之内允许投票
	到期之后删除那个 KeyPostVotedZSetPF
*/

func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.String("userID", strconv.Itoa(int(userID))),
		zap.String("postID", p.PostID),
		zap.Int64("direction", p.Direction))

	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//redis查询id列表
	ids, err := redis.GetCommunityPostIDInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	//mysql查询帖子详情
	//返回的数据还要按照我给定的id数据进行返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//查询每篇帖子的投票信息
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将帖子作者以及分区信息查询出来填充到帖子当中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
			continue
		}

		//根据社区信息查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Int64("CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListNew 按照有无communityID
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
	}
	return
}
