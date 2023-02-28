package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

const secret = "bilibili.com"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户密码错误")
)

// CheckUserExist 根据username判断用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(user_id) from user where username=?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户数据
func InsertUser(user models.User) (err error) {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行sql插入数据库
	sqlStr := "insert into user(user_id,username,password) values(?,?,?)"
	_, err = db.Exec(sqlStr, user.UserId, user.Username, user.Password)
	return
}

// 密码加密
func encryptPassword(Password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(Password)))
}

func Login(user *models.User) (err error) {
	//用户登录的密码
	pwd := user.Password
	//根据username查询用户是否存在
	sqlStr := "select user_id,username,password from user where username=?"
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	password := encryptPassword(pwd)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}
