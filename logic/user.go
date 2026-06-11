package logic

import (
	"GoForum/dao/mysql"
	"GoForum/models/param"
)

func Register(params *param.RegisterParams) error {
	return mysql.InsertUser(params)
}

func Login(params *param.LoginParams) error {
	return mysql.UserLogin(params)
}
