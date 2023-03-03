package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := "insert into post(post_id,title,content,author_id,community_id) values (?,?,?,?,?)"
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostById 根据id查询单个帖子数据
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := "select post_id,title,content,author_id,community_id,create_time from post where post_id=?"
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(pageNum int64, PageSize int64) (posts []*models.Post, err error) {
	posts = make([]*models.Post, 0, PageSize)
	strSql := "select post_id,title,content,author_id,community_id,create_time from post limit ?,?"
	err = db.Select(&posts, strSql, (pageNum-1)*PageSize, PageSize)
	return
}

func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	postList = make([]*models.Post, 0, len(ids))
	strSql := "select post_id,title,content,author_id,community_id,create_time from post where post_id in (?) order by FIND_IN_SET(post_id,?)"
	query, args, err := sqlx.In(strSql, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	db.Select(&postList, query, args...)
	return
}
