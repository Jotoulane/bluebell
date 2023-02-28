package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	//生成UID
	uerId := snowflake.GenID()
	user := models.User{
		UserId:   uerId,
		Username: p.Username,
		Password: p.Password,
	}
	//保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针，就能拿到user.UserID
	if err := mysql.Login(user); err != nil {
		return "", err
	}
	//生成JWT
	return jwt.GenToken(user.UserId, user.Username)
}
