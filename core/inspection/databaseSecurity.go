// Package inspection 安全检查
package inspection

import (
	"dbcheck/core/dbmsg"
	"dbcheck/core/result"
	"dbcheck/pkg/glog"
	"fmt"
	"strings"
)

func BaselineCheckPortDesign() {
	if vi, okk := SecurityCanCheck["databaseport"]; okk {
		var d = make(map[string]string)
		cc := dbmsg.GlobalVariablesMap()
		d["checkStatus"] = "normal" // 正常
		if strings.EqualFold(cc["port"], vi) {
			d["checkStatus"] = "abnormal" // 异常
			d["checkType"] = "databasePort"
			d["threshold"] = "默认端口"
			d["currentValue"] = fmt.Sprintf("port=%s", cc["port"])
			glog.Error(fmt.Sprintf(" US2-01 The MySQL service uses the default port. The information is as follows: using port: %s.", cc["port"]))
		}
		result.R.Security.PortDesign.DatabasePort = append(result.R.Security.PortDesign.DatabasePort, d)
	}
}
