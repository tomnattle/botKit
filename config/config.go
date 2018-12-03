package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

var (
	config *Config
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
	AITS            *AITSConfig   `toml:"AITS"`
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
	WinMode   string `toml:"WinMode"`
	Tsketch   string `toml:"Tsketch"`
	Admin     string `toml:"Admin"`
	Dashboard string `toml:"Dashboard"`
	Center    string `toml:"Center"`
	Interview string `toml:"Interview"`
	Login     string `toml:"Login"`
}

type AITSConfig struct {
	URL string `toml:"URL"`
}

func init() {
	var (
		cfg string
	)
	if cfg = os.Getenv("CFG"); len(cfg) > 0 {
		initConfig(cfg)
		return
	}
	flag.StringVar(&cfg, "cfg",
		"unset config file",
		"configFile is required.\n-cfg ./config/*/config.toml\ndev test prod\n")
	flag.Parse()
	initConfig(cfg)
	return
}

func GetConfig() *Config {
	if config == nil {
		panic("nil config error")
	}
	return config
}

func initConfig(configFile string) {
	config = &Config{}
	_, err := toml.DecodeFile(configFile, config)
	if err != nil {
		panic(fmt.Errorf("init config error, try decode config file %v, %v",
			configFile, err))
	}

	switch config.Environment {
	case "dev", "test", "prod":
	default:
		panic(fmt.Errorf("init config error, try find environment error, %v",
			config.Environment))
	}

	fmt.Printf("init config success, config file path is %s\n", configFile)
}
