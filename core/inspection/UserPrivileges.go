// Package inspection 安全检查
package inspection

import (
	"dbcheck/core/dbmsg"
	"dbcheck/core/result"
	"dbcheck/pkg/db"
	"dbcheck/pkg/glog"
	"fmt"
	"strings"
)

var SecurityCanCheck = map[string]string{
	`anonymousUsers`:                  ``,     // 判断是否存在匿名用户
	`emptyPasswordUser`:               ``,     // 判断是否存在空密码用户
	`rootUserRemoteLogin`:             ``,     // 判断是否存在root用户远程访问
	`normalUserConnectionUnlimited`:   ``,     // 判断是否存在用户连接无限制 user@'%'
	`userPasswordSame`:                ``,     // 判断是否存在多个用户密码相同
	`normalUserDatabaseAllPrivilages`: ``,     // 判断是否存在多个用户密码相同
	`normalUserSuperPrivilages`:       ``,     // 判断是否存在多个用户密码相同
	`databasePort`:                    `3306`, // 判断数据库是否使用默认端口  检测端口的阈值，检测数据库使用端口是否使用3306
}

// UserPrivileges 用户权限检查
func UserPrivileges() {
	var tmpPassword, tmpUser, tmpHost interface{}
	for _, v := range dbmsg.MysqlUser() {
		var d = make(map[string]string)
		d["user"] = v.User
		d["host"] = v.Host
		// 检查匿名用户
		if _, ok := SecurityCanCheck["anonymoususers"]; ok {
			if v.User == "" {
				d["checkStatus"] = "abnormal" // 异常
				d["checkType"] = "anonymousUsers"
				d["threshold"] = "匿名用户"
				d["currentValue"] = fmt.Sprintf("%s@%s", v.User, v.Host)
				glog.Error(fmt.Sprintf(" US1-01 Anonymous users currently exist. The information is as follows: user: %s host: %s", v.User, v.Host))
			} else {
				d["checkStatus"] = "normal" // 正常
				d["checkType"] = "anonymousUsers"
			}
			result.R.Security.UserPriDesign.AnonymousUsers = append(result.R.Security.UserPriDesign.AnonymousUsers, d)
		}
		m := newMap(d)
		// 检查空密码用户
		if _, ok := SecurityCanCheck["emptypassworduser"]; ok {
			if v.Pwd == "" {
				m["checkStatus"] = "abnormal" // 异常
				m["checkType"] = "emptyPasswordUser"
				m["threshold"] = "空密码用户"
				m["currentValue"] = fmt.Sprintf("%s@%s", v.User, v.Host)
				glog.Error(fmt.Sprintf(" US1-02 The current username password is empty. The information is as follows: user: %s host: %s", v.User, v.Host))
			} else {
				m["checkStatus"] = "normal" // 异常
				m["checkType"] = "emptyPasswordUser"
			}
			result.R.Security.UserPriDesign.EmptyPasswordUser = append(result.R.Security.UserPriDesign.EmptyPasswordUser, m)
		}
		n := newMap(d)
		// 检查root用户远端登录，只能本地连接
		if _, ok := SecurityCanCheck["rootuserremotelogin"]; ok {
			if v.User == "root" && v.Host != "localhost" && v.Host != "127.0.0.1" {
				n["checkStatus"] = "abnormal" // 异常
				n["checkType"] = "rootUserRemoteLogin"
				n["threshold"] = "root远端访问"
				n["currentValue"] = fmt.Sprintf("%s@%s", v.User, v.Host)
				glog.Error(fmt.Sprintf(" US1-03 The root user is currently in remote login danger. The information is as follows: user: %s host: %s", v.User, v.Host))
			} else {
				n["checkStatus"] = "normal" // 异常
				n["checkType"] = "rootUserRemoteLogin"
			}
			result.R.Security.UserPriDesign.RootUserRemoteLogin = append(result.R.Security.UserPriDesign.RootUserRemoteLogin, n)
		}
		o := newMap(d)
		// 检查普通用户远端连接的限制，不允许使用%
		if _, ok := SecurityCanCheck["normaluserconnectionunlimited"]; ok {
			if v.User != "" && v.User != "root" && v.Pwd != "" && v.Host == "%" {
				o["checkStatus"] = "abnormal" // 异常
				o["checkType"] = "normalUserConnectionUnlimited"
				o["threshold"] = "普通用户@%"
				o["currentValue"] = fmt.Sprintf("%s@%s", v.User, v.Host)
				glog.Error(fmt.Sprintf(" US1-04 The current user name has no connection IP limit. The information is as follows: user: %s host: %s", v.User, v.Host))
			} else {
				o["checkStatus"] = "normal" // 异常
				o["checkType"] = "normalUserConnectionUnlimited"
			}
			result.R.Security.UserPriDesign.NormalUserConnectionUnlimited = append(result.R.Security.UserPriDesign.NormalUserConnectionUnlimited, o)
		}
		// 检查不同用户使用相同密码
		p := newMap(d)
		if _, ok := SecurityCanCheck["userpasswordsame"]; ok {
			if v.Pwd == tmpPassword {
				p["checkStatus"] = "abnormal" // 异常
				p["checkType"] = "userPasswordSame"
				p["threshold"] = "密码相同用户"
				p["currentValue"] = fmt.Sprintf("%s@%s", v.User, v.Host)
				glog.Error(fmt.Sprintf(" US1-05 Different users in the current database use the same password, please change it. The information is as follows: user1: %v@%v  user2: %s@%s", tmpUser, tmpHost, v.User, v.Host))
			} else {
				p["checkStatus"] = "normal" // 异常
				p["checkType"] = "userPasswordSame"
			}
			result.R.Security.UserPriDesign.UserPasswordSame = append(result.R.Security.UserPriDesign.UserPasswordSame, p)
		}
		tmpPassword = v.Pwd
		tmpUser = v.User
		tmpHost = v.Host

		// 检查跨用户权限*.*
		strSql := fmt.Sprintf("show grants for '%s'@'%s'", v.User, v.Host)
		var cd string
		db.Get(&cd, strSql)
		if v.User != "root" && v.Host != "localhost" && v.Host != "127.0.0.1" {
			// 检查当前用户是否存在ON *.*
			q := newMap(d)
			if _, ok := SecurityCanCheck["normaluserdatabaseallprivilages"]; ok {
				if strings.Contains(cd, "ON *.*") {
					q["checkStatus"] = "abnormal" // 异常
					q["checkType"] = "normalUserDatabaseAllPrivilages"
					q["threshold"] = "普通用户 ON *.*"
					q["currentValue"] = fmt.Sprintf("%s@%s", v.User, v.Host)
					glog.Error(fmt.Sprintf(" US1-06 Cross-user permissions currently exist (ON *.*). The information is as follows: user@host: %s@%s", v.User, v.Host))
				} else {
					q["checkStatus"] = "normal" // 异常
					q["checkType"] = "normalUserDatabaseAllPrivilages"
				}
				result.R.Security.UserPriDesign.NormalUserDatabaseAllPrivilages = append(result.R.Security.UserPriDesign.NormalUserDatabaseAllPrivilages, q)
			}
			r := newMap(d)
			// 检查当前用户是否WITH GRANT OPTION
			if _, ok := SecurityCanCheck["normalusersuperprivilages"]; ok {
				if strings.Contains(cd, "WITH GRANT OPTION") {
					r["checkStatus"] = "abnormal" // 异常
					r["checkType"] = "normalUserSuperPrivilages"
					r["threshold"] = "普通用户super权限"
					r["currentValue"] = fmt.Sprintf("%s@%s", v.User, v.Host)
					glog.Error(fmt.Sprintf(" US1-07 The current user has permission transfer (WITH GRANT OPTION). The information is as follows: user@host: %s@%s", v.User, v.Host))
				} else {
					r["checkStatus"] = "normal" // 正常
					r["checkType"] = "normalUserSuperPrivilages"
				}
				result.R.Security.UserPriDesign.NormalUserSuperPrivilages = append(result.R.Security.UserPriDesign.NormalUserSuperPrivilages, r)
			}
		}
	}
}
