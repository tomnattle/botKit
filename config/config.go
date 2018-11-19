package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
)

var (
	configFile string
	config     *Config
)

type Config struct {
	AppID           int           `toml:"AppID"`
	SubAppID        int           `toml:"SubAppID"`
	AppName         string        `toml:"AppName"`
	ProductID       int           `toml:"ProductID"`
	Environment     string        `toml:"Environment"`
	Addr            string        `toml:"Addr"`
	InsideSignature string        `toml:"InsideSignature"`
	Logger          *LoggerConfig `toml:"Log"`
	MySQL           *MySQLConfig  `toml:"MySQL"`
	Redis           *RedisConfig  `toml:"Redis"`
	SMS             *SMSConfig    `toml:"SMS"`
	URI             *URIConfig    `toml:"URI"`
	Dfs             *DfsConfig    `toml:"Dfs"`
}

type DfsConfig struct {
	Server string `toml:"Server"`
}

type LoggerConfig struct {
	LogFile string `toml:"LogFile"`
}

type MySQLConfig struct {
	Username string `toml:"UserName"`
	Password string `toml:"Password"`
	Addr     string `toml:"Addr"`
	DB       string `toml:"DB"`
}

type RedisConfig struct {
	IsCluster bool   `toml:"IsCluster"`
	Addr      string `toml:"Addr"`
	PoolSize  int    `toml:"PoolSize"`
	Prefix    string `toml:"Prefix"`
}

type SMSConfig struct {
	URL                     string `toml:"URL"`
	Account                 string `toml:"Account"`
	Model                   string `toml:"Model"`
	Secret                  string `toml:"Secret"`
	DurationMinutes         int    `toml:"DurationMinutes"`
	AuthCodeDurationMinutes int    `toml:"AuthCodeDurationMinutes"`
}

type URIConfig struct {
	ChatBot   string `toml:"ChatBot"`
	WinMode   string `toml:"WinMode"`
	Tsketch   string `toml:"Tsketch"`
	Admin     string `toml:"Admin"`
	Dashboard string `toml:"Dashboard"`
}

func init() {
	flag.StringVar(&configFile, "cfg",
		"unset config file",
		"configFile is required.\n-cfg ./config/*/config.toml\ndev test prod\n")
	flag.Parse()
	initConfig()
	initEnvironment()
}

func GetConfig() *Config {
	if config == nil {
		panic("nil config error")
	}
	return config
}

func initEnvironment() {
	switch config.Environment {
	case "dev", "test", "prod":
	default:
		panic(fmt.Errorf("init config error, try find environment error, %v",
			config.Environment))
	}
}

func initConfig() {
	config = &Config{}
	_, err := toml.DecodeFile(configFile, config)
	if err != nil {
		panic(fmt.Errorf("init config error, try decode config file %v, %v",
			configFile, err))
	}
}
