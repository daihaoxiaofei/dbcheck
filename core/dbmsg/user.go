package dbmsg

import (
	"dbcheck/pkg/db"
	"sync"
)

type mysqlUser struct {
	User string `db:"user"`
	Host string `db:"host"`
	Pwd  string `db:"authentication_string"`
}

var (
	user     []mysqlUser // 用户权限 mysql.user
	userOnce sync.Once   // 用户权限 mysql.user
)

func MysqlUser() []mysqlUser {
	userOnce.Do(func() {
		db.Select(&user, "select user,host,authentication_string from mysql.user")
	})
	return user
}
