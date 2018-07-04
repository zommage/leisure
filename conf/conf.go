package conf

import (
	ini "gopkg.in/ini.v1"
)

type Config struct {
	BaseConf
	DbConf
	CollectGrpc
	LogConf
}

type BaseConf struct {
	HttpPort   string `ini:"HttpPort"`   // http port
	HttpsPort  string `ini:"HttpsPort"`  // https port
	Env        string `ini:"Env"`        // 环境信息
	Secret     string `ini:"Secret"`     // 签名秘钥
	RsaSertKey string `ini:"RsaSertKey"` // rsa 私钥地址
	RsaPubKey  string `ini:"RsaPubKey"`  // rsa 公钥地址
	SslKey     string `ini:"SslKey"`     // ssl key 地址
	SslCrt     string `ini:"SslCrt"`     // ssl crt 地址
}

// mysql db config
type DbConf struct {
	DbName        string `ini:"DbName"`
	DbHost        string `ini:"DbHost"`
	DbPort        string `ini:"DbPort"`
	DbUser        string `ini:"DbUser"`
	DbPassword    string `ini:"DbPassword"`
	DbLogEnable   bool   `ini:"DbLogEnable"`
	DbMaxConnect  int    `ini:"DbMaxConnect"`
	DbIdleConnect int    `ini:"DbIdleConnect"`
}

type CollectGrpc struct {
	RoomGrpc string `ini:"RoomGrpc"` //房间服务
}

// Log config
type LogConf struct {
	LogPath  string `ini:"LogPath"`
	LogLevel string `ini:"LogLevel"`
}

var Conf *Config

func InitConfig(confPath *string) (*Config, error) {
	Conf = new(Config)
	if err := ini.MapTo(Conf, *confPath); err != nil {
		return nil, err
	}

	return Conf, nil
}
