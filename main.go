package main

import (
	"dbcheck/core/output/pdf"
	"time"

	"dbcheck/core/inspection"
	"dbcheck/core/result"
	"dbcheck/pkg/glog"
)

func main() {
	// 配置文件初始化
	result.BeginTime = time.Now()

	inspection.Check() // 巡检

	result.ConsumingTime = time.Since(result.BeginTime).Seconds()
	pdf.OutPdf()

	glog.Info(`完成`)
}
