package alertinit

import (
	"fmt"
	"os"
	"xxl_job_alert/alertmodel"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var Conf alertmodel.Config

// 请勿修改这里的文件地址，是为了保持windows开发环境与docker容器内运行环境配置文件路径一致
const configPath = "./etc/xxl_job_alert/conf"
const logPath = "./etc/xxl_job_alert/log"

func Init() {
	initZeroLog()
	loadConfToml()
}

func loadConfToml() {
	// 初始化 Viper 配置
	viper.SetConfigType("toml")
	viper.SetConfigName("conf")     // 配置文件名，不包含扩展名
	viper.AddConfigPath(configPath) // 配置文件路径，此处假设配置文件与程序在同一目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	// 解析配置到结构体
	if err := viper.Unmarshal(&Conf); err != nil {
		fmt.Println("Error unmarshaling config:", err)
		return
	}
}
func initZeroLog() {
	fmt.Println("加载日志文件开始")
	logFilePath := fmt.Sprintf("%s/log.log", logPath)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.MkdirAll(logPath, 0666)
		os.Create(logFilePath)
	}
	logfile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Error().Msgf("日志文件错误: %v", err)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(logfile).With().
		Caller().
		Timestamp().
		Logger()
	// defer logfile.Close()
	fmt.Println("加载日志文件结束")
}
