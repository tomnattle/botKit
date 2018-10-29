package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
)

var (
	environment Environment
	configFile  string
	config      *Config
)

type Config struct {
	AppID       int              `toml:"AppID"`
	SubAppID    int              `toml:"SubAppID"`
	AppName     string           `toml:"AppName"`
	Environment string           `toml:"Environment"`
	Addr        string           `toml:"Addr"`
	Logger      *LoggerConfig    `toml:"Log"`
	MySQL       *MySQLConfig     `toml:"MySQL"`
	Redis       *RedisConfig     `toml:"Redis"`
	SMS         *SMSConfig       `toml:"SMS"`
	Signature   *SignatureConfig `toml:"Signature"`
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

type SignatureConfig struct {
	SecretKey string `toml:"SecretKey"`
}

func init() {
	flag.StringVar(&configFile, "cfg",
		"unset config file",
		"configFile is required.\n-cfg ./config/*/config.toml\ndev test prod\n")
	flag.Parse()
	initConfig()
	initEnvironment()
}

type Environment uint8

const (
	_ Environment = iota
	DEV
	TEST
	PROD
)

func (e Environment) String() string {
	switch e {
	case DEV:
		return "dev"
	case TEST:
		return "test"
	case PROD:
		return "prod"
	}
	return "runtime find environment error !!!"
}

func GetConfig() *Config {
	if config == nil {
		panic("nil config error")
	}
	return config
}

func GetEnvironment() Environment {
	return environment
}

func initEnvironment() {
	switch config.Environment {
	case "dev":
		environment = DEV
	case "test":
		environment = TEST
	case "prod":
		environment = PROD
	}
	if uint8(environment) == 0 {
		panic(fmt.Errorf("init config error, try find environment error"))
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
