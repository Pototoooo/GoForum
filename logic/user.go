package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models/param"
)

func Register(params *param.RegisterParams) error {
	return mysql.InsertUser(params)
}

func Login(params *param.LoginParams) error {
	return mysql.UserLogin(params)
}
