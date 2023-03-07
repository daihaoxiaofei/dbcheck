package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"dbcheck/pkg/config"
	"dbcheck/pkg/glog"
)

var DB *sqlx.DB

// InitPgDB Init postgres Db
func init() {
	Cf := config.C.DBInfo
	source := fmt.Sprintf("root:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms",
		Cf.Password, Cf.Host, Cf.Database)
	var err error
	DB, err = sqlx.Connect("mysql", source) // "user=postgres dbname=order password=a123456 host=192.168.20.74 sslmode=disable")
	if err != nil {
		panic(`sqlx连接到从数据库出现错误: ` + err.Error())
	}
	// 设置连接池中的最大连接数。
	DB.SetMaxOpenConns(800)

	// 设置最大的空闲连接数  设置小于等于0的数意味着不保留空闲连接。
	DB.SetMaxIdleConns(15)

	DB.SetConnMaxLifetime(time.Hour) // 连接过期时间 如不设置 连接会被一直保持
}

// GetMap 查询数据库，结果集为两列，切第一列的数据是第二列的key，对列数据生成key-value并返回（不包含列名）
func GetMap(sqlStr string) map[string]string {
	rows, err := DB.Query(sqlStr)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	result := make(map[string]string)
	var k, v string
	for rows.Next() {
		rows.Scan(&k, &v)
		result[k] = v
	}

	return result
}

func Select(dest interface{}, query string, args ...interface{}) {
	err := DB.Select(dest, query, args...)
	if err != nil {
		glog.Error(`sql.Select`, zap.Error(err))
	}
}

func Get(dest interface{}, query string, args ...interface{}) {
	err := DB.Get(dest, query, args...)
	if err != nil {
		glog.Error(`sql.Get`, zap.Error(err))
	}
}

func Exec(query string, args ...any) sql.Result {
	Result, err := DB.Exec(query, args...)
	if err != nil {
		glog.Error(`sql.Exec`, zap.Error(err), zap.String(`query`, query), zap.Any(`args`, args))
	}
	return Result
}
