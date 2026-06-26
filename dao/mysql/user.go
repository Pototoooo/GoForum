package mysql

import (
	"errors"
	"fmt"

	"GoForum/models/param"
	"GoForum/pkg/snowflake"
	"GoForum/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"go.uber.org/zap"
)

var (
	ErrorUserExsist   = errors.New("用户已存在")
	ErrorUserNotFound = errors.New("用户不存在")
)
var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		settings.Config.Mysql.User,
		settings.Config.Mysql.Password,
		settings.Config.Mysql.Host,
		settings.Config.Mysql.Port,
		settings.Config.Mysql.DBName,
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("Connect failed", zap.Error(err))
		return
	}
	if err = ensurePostVoteTable(); err != nil {
		zap.L().Error("ensure post_vote table failed", zap.Error(err))
		return
	}
	db.SetMaxIdleConns(settings.Config.Mysql.MaxIdleConns)
	db.SetMaxOpenConns(settings.Config.Mysql.MaxOpenConns)
	return
}

func Close() {
	db.Close()
}

// GetDB 返回数据库连接（用于测试）
func GetDB() *sqlx.DB {
	return db
}

// 数据库查找用户id
func QueryUserIdByName(name string) int64 {
	sqlStr := "select user_id from user where username = ?"
	var id int64
	err := db.Get(&id, sqlStr, name)
	if err != nil {
		return 0
	}
	return id
}

// 根据id查用户信息
func GetUserNameByID(id int64) (username string, err error) {
	sqlStr := "select username from user where user_id = ?"
	err = db.Get(&username, sqlStr, id)
	if err != nil {
		return "", err
	}
	return username, nil
}

// 数据库注册新用户
func InsertUser(p *param.RegisterParams) (err error) {
	// 数据库查找是否存在
	temp := QueryUserIdByName(p.Username)
	if temp != 0 {
		return ErrorUserExsist
	}
	// 生成id
	id := snowflake.GenerateID()
	// 加密密码
	hashedPassworded := HashPassword(p.Password)
	// sql添加新用户
	sqlStr := "insert into user(user_id,username,password) values (?,?,?)"
	_, err = db.Exec(sqlStr, id, p.Username, hashedPassworded)
	if err != nil {
		return err
	}
	return
}

// 用户登录
func UserLogin(p *param.LoginParams) (err error) {
	// 检查用户是否存在
	id := QueryUserIdByName(p.Username)
	if id == 0 {
		return ErrorUserNotFound
	}
	// 哈希化并检查密码
	sqlStr := "select password from user where username = ?"
	var dbPassword string
	err = db.Get(&dbPassword, sqlStr, p.Username)
	if err != nil {
		return errors.New("密码错误")
	}
	// 对比密码
	if !CheckPassword(dbPassword, p.Password) {
		return errors.New("密码错误")
	}
	return
}

// 加密密码
func HashPassword(password string) string {
	// 使用bcrypt进行密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

// 验证密码
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
