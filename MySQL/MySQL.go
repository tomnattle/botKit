package MySQL

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ifchange/botKit/config"
	"time"
)

var (
	cfg                 *config.MySQLConfig
	DefaultConnNotClose *sql.DB
)

func init() {
	cfg = config.GetConfig().MySQL
	if cfg == nil {
		panic("MySQL config is nil")
	}

	conn, err := getConn(cfg)
	if err != nil {
		panic(fmt.Errorf("MySQL get connection error %v %v", cfg.Addr, err))
	}
	err = conn.Ping()
	if err != nil {
		panic(fmt.Errorf("MySQL connection error %v", err))
	}

	DefaultConnNotClose = conn
	DefaultConnNotClose.SetConnMaxLifetime(time.Duration(3) * time.Second)
	DefaultConnNotClose.SetMaxOpenConns(1000)
	DefaultConnNotClose.SetMaxIdleConns(100)
}

func GetConn() (*sql.DB, error) {
	return getConn(cfg)
}

func getConn(cfg *config.MySQLConfig) (*sql.DB, error) {
	// "用户名:密码@tcp(IP:端口)/数据库?charset=utf8mb4&sql_notes=false&sql_notes=false&timeout=90s&collation=utf8mb4_general_ci&parseTime=True&loc=Local"
	source := getConnStr(cfg)
	return sql.Open("mysql", source)
}

func getConnStr(cfg *config.MySQLConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&sql_notes=false&sql_notes=false&timeout=90s&collation=utf8mb4_general_ci&parseTime=True&loc=Local",
		cfg.Username, cfg.Password,
		cfg.Addr, cfg.DB)
}

func GetConnStr(cfg *config.MySQLConfig) string {
	return getConnStr(cfg)
}
