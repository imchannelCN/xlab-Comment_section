package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func init() {
	// viper.SetConfigName("config1") // 读取yaml配置文件
	viper.SetConfigName("config") // 读取json配置文件
	viper.AddConfigPath("/")      //设置配置文件的搜索目录
	//viper.AddConfigPath("$HOME/.appname")  // 设置配置文件的搜索目录
	viper.AddConfigPath(".") // 设置配置文件和可执行二进制文件在用一个目录
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Println("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}

	viper_main()
}

func viper_main() {
	fmt.Println("获取配置文件的port", viper.GetInt("port"))
	fmt.Println("获取配置文件的mysql.url", viper.GetString(`mysql.url`))
	fmt.Println("获取配置文件的mysql.username", viper.GetString(`mysql.username`))
	fmt.Println("获取配置文件的mysql.password", viper.GetString(`mysql.password`))
	fmt.Println("over")
}
