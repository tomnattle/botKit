package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"strings"
)

func LoadConfig() error {
	return nil
}

type Config struct {
	AppID     int              `toml:"AppID"`
	SubAppID  int              `toml:"SubAppID"`
	AppName   string           `toml:"AppName"`
	Addr      string           `toml:"Addr"`
	Logger    *LoggerConfig    `toml:"Log"`
	MySQL     *MySQLConfig     `toml:"MySQL"`
	Redis     *RedisConfig     `toml:"Redis"`
	SMS       *SMSConfig       `toml:"SMS"`
	Signature *SignatureConfig `toml:"Signature"`
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
	URL             string `toml:"URL"`
	Account         string `toml:"Account"`
	Model           string `toml:"Model"`
	Secret          string `toml:"Secret"`
	DurationMinutes int    `toml:"DurationMinutes"`
}

type SignatureConfig struct {
	SecretKey string `toml:"SecretKey"`
}

func init() {
	initEnvironment()
	initConfig(environment.configFile())
}

type Environment uint8

const (
	_ Environment = iota
	DEV
	TEST
	PROD
)

func (e *Environment) String() string {
	switch *e {
	case DEV:
		return "dev"
	case TEST:
		return "test"
	case PROD:
		return "prod"
	}
	return "runtime find environment error !!!"
}

func (e *Environment) configFile() (configFile string) {
	if e == nil {
		panic("un find environment, start fail")
	}
	runAbs, err := filepath.Abs(".")
	if err != nil {
		panic("try to find run abs path error")
	}
	pathS := strings.Split(runAbs, AppName)
	rootParentPath := pathS[0]
	env := ""
	switch *e {
	case DEV:
		env = "dev"
	case TEST:
		env = "test"
	case PROD:
		env = "prod"
	}
	configFile = filepath.Join(rootParentPath, AppName,
		"config", env, "config.toml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file not exist in %s", configFile))
	}
	return
}

var (
	environment *Environment
	config      *Config
)

func GetConfig() *Config {
	if config == nil {
		panic("nil config error")
	}
	return config
}

func GetEnvironment() *Environment {
	return environment
}

func initEnvironment() {
	parser :=
		func(e string) {
			switch e {
			case "dev":
				formatE := DEV
				environment = &formatE
			case "test":
				formatE := TEST
				environment = &formatE
			case "prod":
				formatE := PROD
				environment = &formatE
			}
		}

	if environment == nil {
		parser(os.Getenv("ENVIRONMENT"))
	}

	if environment == nil {
		parser(*(flag.String("e", "dev",
			"set env, e.g dev test prod")))
	}
}

func initConfig(configFile string) {
	config = &Config{}
	_, err := toml.DecodeFile(configFile, config)
	if err != nil {
		panic(fmt.Errorf("init config error, try decode config file %v, %v",
			configFile, err))
	}
}
