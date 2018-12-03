// Package ormx wraps ORM and provide some basic operation.
package xorm

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/ifchange/botKit/config"
	"time"
)

// list of DB errors
var (
	ErrNotExist = errors.New("not exist")
)

var (
	cfg *config.MySQLConfig
	orm *xorm.Engine
)

// Init inits db info and setting db.
func init() {
	cfg = config.GetConfig().MySQL
	if cfg == nil {
		panic("MySQL config is nil")
	}

	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&sql_notes=false&sql_notes=false&timeout=90s&collation=utf8mb4_general_ci&parseTime=True&loc=Local",
		cfg.Username, cfg.Password,
		cfg.Addr, cfg.DB)

	var err error
	orm, err = xorm.NewEngine("mysql", url)
	if err != nil {
		panic(fmt.Errorf("MySQL get connection error %v %v", cfg.Addr, err))
	}

	orm.SetConnMaxLifetime(time.Duration(3) * time.Second)
	orm.SetMaxOpenConns(1000)
	orm.SetMaxIdleConns(100)
}

// ORM returns initialized orm.
func ORM() *xorm.Engine {
	return orm
}

// GetByID return a obj by id.
func GetByID(id int64, obj interface{}) error {
	has, err := orm.Id(id).Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return ErrNotExist
	}
	return nil
}

// SoftDeleteByID set record's is_deleted to true.
func SoftDeleteByID(id int, tab Table) error {
	tab.SetDeleted()
	_, err := orm.Id(id).Update(tab)
	return err
}

// DeleteByID delete a record from database.
func DeleteByID(id int, obj interface{}) error {
	_, err := orm.Id(id).Delete(obj)
	return err
}

// Create insert a record into database.
func Create(obj interface{}) error {
	_, err := orm.Insert(obj)
	return err
}
