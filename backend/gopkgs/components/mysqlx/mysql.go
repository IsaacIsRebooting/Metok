package mysqlx

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/components"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	globalClientMap = sync.Map{}
	globalConfigMap = make(components.ConfigMap[*Config])
)

func GetConfig() components.ConfigMap[*Config] {
	return globalConfigMap
}

func Init(cm components.ConfigMap[*Config]) (func() error, error) {
	globalConfigMap = cm
	for k, v := range cm {
		db, err := Connect(v)
		if err != nil {
			return nil, err
		}
		globalClientMap.Store(k, db)
	}
	return IsHealth, nil
}

// Connect 使用给定的配置信息连接到数据库，并返回一个 gorm.DB 实例。
// 这个函数首先会设置配置的默认值，然后根据配置信息初始化一个 SQL 数据库连接。
// 它还会根据配置设置连接池的相关参数，以优化数据库的连接和使用。
// 参数:
//
//	c *Config - 数据库配置信息的指针。
//
// 返回值:
//
//	*gorm.DB - 初始化后的 GORM 数据库实例指针。
//	error - 如果在连接数据库过程中遇到任何错误，都会返回该错误。
func Connect(c *Config) (*gorm.DB, error) {
	// 设置配置的默认值，确保所有必要的配置项都有有效的值。
	c.setDefault()

	// 根据配置信息中的数据源名称 (DSN) 初始化一个 SQL 数据库连接。
	originDB, err := sql.Open("mysql", c.ToDSN())
	if err != nil {
		return nil, err
	}

	// 设置连接池的最大空闲连接数。
	originDB.SetMaxIdleConns(c.MaxIdle)
	// 设置连接池的最大打开连接数。
	originDB.SetMaxOpenConns(c.MaxOpen)

	// 根据配置设置连接池中连接的最大空闲时间。
	if c.ConnMaxIdleTime > 0 {
		originDB.SetConnMaxIdleTime(time.Duration(c.ConnMaxIdleTime) * time.Second)
	} else {
		originDB.SetConnMaxIdleTime(0)
	}

	// 根据配置设置连接池中连接的最大生命周期。
	if c.ConnMaxLifeTime > 0 {
		originDB.SetConnMaxLifetime(time.Duration(c.ConnMaxLifeTime) * time.Second)
	} else {
		originDB.SetConnMaxLifetime(0)
	}

	// 将初始化的数据库连接设置为 GORM 的连接池。
	var connPoll gorm.ConnPool = originDB

	// 使用 GORM 的 MySQL 驱动配置一个新的 Dialector。
	dialector := mysql.New(mysql.Config{
		Conn: connPoll,
	})

	// 使用配置好的 Dialector 初始化 GORM，并返回 GORM 的数据库实例。
	// 设置命名策略为使用单数表名。
	return gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
}

func GetDB(ctx context.Context, keys ...string) *gorm.DB {
	key := getKey(keys...)

	value, ok := globalClientMap.Load(key)
	if !ok {
		panic(fmt.Sprintf("%s not init", key))
	}
	return value.(*gorm.DB)

}
func getKey(keys ...string) string {
	if len(keys) == 0 {
		return "default"
	}
	return keys[0]
}
func IsHealth() (err error) {
	globalClientMap.Range(func(key, value interface{}) bool {
		client := value.(*gorm.DB)
		db, e := client.DB()
		if e != nil {
			err = e
			return false
		}

		err = db.Ping()
		if err != nil {
			return false
		}

		log.Infof("mysql %s health check ok", key)
		return true
	})
	return err
}
