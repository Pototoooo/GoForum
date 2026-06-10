package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	//"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Config *config

type config struct {
	Port      int          `mapstructure:"port"`      // 端口号
	StartTime string       `mapstructure:"starttime"` // 雪花算法起始时间
	MachineID int64        `mapstructure:"machineid"` // 雪花算法机器ID
	Log       *LogConfig   `mapstructure:"log"`       // 日志级别
	Redis     *RedisConfig `mapstructure:"redis"`     // redis配置
	Mysql     *MysqlConfig `mapstructure:"mysql"`     // mysql配置
}
type LogConfig struct {
	Level      string `mapstructure:"level"`      // 日志级别
	Filename   string `mapstructure:"filename"`   // 日志文件名
	MaxSize    int    `mapstructure:"maxsize"`    // 日志文件最大大小（MB）
	MaxBackups int    `mapstructure:"maxbackups"` // 日志文件最大备份数量
	MaxAge     int    `mapstructure:"maxage"`     // 日志文件最大年龄（天）
	Console    bool   `mapstructure:"console"`    // 是否输出到日志
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`     // redis主机
	Port     int    `mapstructure:"port"`     // redis端口
	Password string `mapstructure:"password"` // redis密码
	DB       int    `mapstructure:"db"`       // redis数据库索引
	PoolSize int    `mapstructure:"poolsize"` // redis连接池大小
}
type MysqlConfig struct {
	Host         string `mapstructure:"host"`         // mysql主机
	Port         int    `mapstructure:"port"`         // mysql端口
	User         string `mapstructure:"user"`         // mysql用户名
	Password     string `mapstructure:"password"`     // mysql密码
	DBName       string `mapstructure:"dbname"`       // mysql数据库名
	MaxOpenConns int    `mapstructure:"maxopenconns"` // mysql最大打开连接数
	MaxIdleConns int    `mapstructure:"maxidleconns"` // mysql最大空闲连接数Conn3秒
}

func Init() (err error) {
	Config = &config{}
	// 1. 设置 Viper 配置信息
	viper.SetConfigFile("config.yaml")
	// viper.SetConfigName("config") // 指定配置文件名称（不需要带后缀）
	// viper.SetConfigType("yaml")   // 指定配置文件类型
	viper.AddConfigPath(".") // 指定查找配置文件的路径（这里使用相对路径）

	// 2. 读取配置信息
	// 注意：图片中用了 err := ，在有命名返回值的函数里这会产生变量遮蔽（Shadowing）
	err = viper.ReadInConfig()
	if err != nil { // 读取配置信息失败
		fmt.Println(err)
		return
	}
	// 初次初始化config
	err = viper.Unmarshal(&Config)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 3. 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(&Config); err != nil {
			fmt.Println("unmarshal config failed:", err)
		}
		fmt.Println("config file has changed")
	})

	// 3. 初始化 Gin 并运行
	// r := gin.Default()
	// // 通过 viper.Get 获取端口并格式化字符串
	// if err := r.Run(fmt.Sprintf(":%d", viper.Get("port"))); err != nil {
	// 	panic(err)
	// }

	return err
}
