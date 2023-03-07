package inspection

import (
	"dbcheck/core/dbmsg"
	"dbcheck/core/result"
	"dbcheck/core/stream"
	"dbcheck/pkg/glog"
	"fmt"
	"strconv"
	"strings"
)

const (
	MtinyInt                  = 127
	MunsignedTinyInt          = 255
	MsmallInt                 = 32767
	MunsignedSmallInt         = 65535
	MmediumInt                = 8388607
	MunsignedMediumInt        = 16777215
	Mint                      = 2147483647
	MunsignedInt              = 4294967295
	MbigInt                   = 9223372036854775807
	MunsignedBigint    uint64 = 18446744073709551615
)

var PerformanceCanCheck = map[string]string{
	`binlogDiskUsageRate`:              `100`,         // 检测binlog磁盘利用率			    	binlog落盘是使用磁盘的百分比阈值 >100%
	`historyConnectionMaxUsageRate`:    `80`,          // 检测历史连接数最大使用率			    	最大连接数使用率百分比阈值 >80%
	`tmpDiskTableUsageRate`:            `25`,          // 检测临时表磁盘使用率			    	磁盘临时表使用率百分比阈值 >25%
	`tmpDiskfileUsageRate`:             `25`,          // 检测临时磁盘文件使用率			    	磁盘临时文件使用率百分比阈值  >25%
	`innodbBufferPoolUsageRate`:        `80`,          // 检测innodb buffer pool使用率			    	innodb buffer pool利用率百分比阈值，<80%
	`innodbBufferPoolDirtyPagesRate`:   `50`,          // 检测innodb buffer pool 中脏页率			    	脏页率百分比阈值     >50%
	`innodbBufferPoolHitRate`:          `99`,          // 检测innodb buffer pool命中率			    	innodb buffer pool命中率阈值   <99%
	`openFileUsageRate`:                `75`,          // 检测文件句柄使用率			    	文件句柄使用率百分比阈值         >75%
	`openTableCacheUsageRate`:          `80`,          // 检测表缓存使用率			    	表缓存使用率百分比阈值           >80%
	`openTableCacheOverflowsUsageRate`: `10`,          // 检测表缓存溢出率			    	表缓存溢出百分比阈值             >10%
	`selectScanUsageRate`:              `10`,          // 检测查询发生全表扫描率			    	select查询全表扫描百分比阈值      >10%
	`selectfullJoinScanUsageRate`:      `10`,          // 联表查询发生全表扫描率			    	发生联表查询全表扫描百分比阈值      >10%
	`tableAutoPrimaryKeyUsageRate`:     `85`,          // 检测表自增主键int类型使用率			    	int有符号和无符号的使用率百分比阈值  >85%
	`tableRows`:                        `5000000,10`,  // 检查表行数			    	第一个参数为表行数阈值，第二个参数为平均行长度阈值，以逗号分割。 >5000000,10k
	`diskFragmentationRate`:            `6,30`,        // 检查表空间磁盘碎片率			    	第一个参数为单表表空间大小阈值（单位为G），第二个为磁盘碎片率（单位%）阈值，以逗号分割。单表大于6G，磁盘碎片率>30%
	`bigTable`:                         `10000000,30`, // 检测是否存在大表			    	第一个参数为表的行数量阈值，第二个参数为表空间大小阈值（单位为G）以逗号分割。单表大于1千万，表空间大小>30G
	`coldTable`:                        `7`,           // 检测当前表是否为冷表			    	检测当前表7day内没有进行更新的表
}

func DatabasePerformanceStatusCheck() {
	for k, v := range PerformanceCanCheck {
		// 统计使用磁盘的binlog写入占使用内存buffer的binlog写入的百分比，大于100%则需要增加binlog_cache_size
		if strings.EqualFold(k, "binlogDiskUsageRate") {
			binlogDiskUsageRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["Binlog_cache_disk_use"], dbmsg.GlobalStatusMap()["Binlog_cache_use"])
			if tmpint, _ := strconv.Atoi(v); binlogDiskUsageRate > tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "binlogDiskUsageRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "binlogDiskUsageRate", strconv.Itoa(binlogDiskUsageRate))
				result.R.Performance.PerformanceStatus.BinlogDiskUsageRate = append(result.R.Performance.PerformanceStatus.BinlogDiskUsageRate, d)
				glog.Error(" PF1-01 The current database binlog is using too many disk writes. It is recommended to modify the binlog_cache_size parameter")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "binlogDiskUsageRate"
				d["checkStatus"] = "normal" // 正常
				result.R.Performance.PerformanceStatus.BinlogDiskUsageRate = append(result.R.Performance.PerformanceStatus.BinlogDiskUsageRate, d)
			}
		}
		// 统计历史连接数最大使用率，使用创建过
		if strings.EqualFold(k, "historyConnectionMaxUsageRate") {
			historyConnectionMaxUsageRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["Threads_created"], dbmsg.GlobalVariablesMap()["max_connections"])
			if tmpint, _ := strconv.Atoi(v); historyConnectionMaxUsageRate > tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "historyConnectionMaxUsageRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "historyConnectionMaxUsageRate", strconv.Itoa(historyConnectionMaxUsageRate))
				result.R.Performance.PerformanceStatus.HistoryConnectionMaxUsageRate = append(result.R.Performance.PerformanceStatus.HistoryConnectionMaxUsageRate, d)
				glog.Error(" PF1-02 If the maximum usage of historical database connections exceeds 80%, change the max_connections value and check services")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "historyConnectionMaxUsageRate"
				d["checkStatus"] = "normal" // 正常
				d["currentValue"] = fmt.Sprintf("%s=%s", "historyConnectionMaxUsageRate", strconv.Itoa(historyConnectionMaxUsageRate))
				result.R.Performance.PerformanceStatus.HistoryConnectionMaxUsageRate = append(result.R.Performance.PerformanceStatus.HistoryConnectionMaxUsageRate, d)
			}
		}

		// 统计数据库使用中使用磁盘临时表占使用内存临时表的占用比例Created_tmp_disk_tables/Created_tmp_tables *100% <=25%
		if strings.EqualFold(k, "tmpDiskTableUsageRate") {
			tmpDiskTableUsageRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["Created_tmp_disk_tables"], dbmsg.GlobalStatusMap()["Created_tmp_tables"])
			if tmpint, _ := strconv.Atoi(v); tmpDiskTableUsageRate > tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "tmpDiskTableUsageRate"
				d["checkStatus"] = "abnormal" // 异常
				d["binlogDiskUsageRate"] = strconv.Itoa(tmpDiskTableUsageRate)
				d["currentValue"] = fmt.Sprintf("%s=%s", "tmpDiskTableUsageRate", strconv.Itoa(tmpDiskTableUsageRate))
				result.R.Performance.PerformanceStatus.TmpDiskTableUsageRate = append(result.R.Performance.PerformanceStatus.TmpDiskTableUsageRate, d)
				glog.Error(" PF1-03 Too many disk temporary tables are being used. Check the slow SQL log or parameters tmp_table_size and max_heap_table_size")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "tmpDiskTableUsageRate"
				d["checkStatus"] = "normal" // 正常
				d["currentValue"] = fmt.Sprintf("%s=%s", "tmpDiskTableUsageRate", strconv.Itoa(tmpDiskTableUsageRate))
				result.R.Performance.PerformanceStatus.TmpDiskTableUsageRate = append(result.R.Performance.PerformanceStatus.TmpDiskTableUsageRate, d)
			}
		}

		// 统计数据库使用中使用磁盘临时表占使用内存临时表的占用比例Created_tmp_disk_tables/Created_tmp_tables *100% <=25%
		if strings.EqualFold(k, "tmpDiskfileUsageRate") {
			tmpDiskfileUsageRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["Created_tmp_files"], dbmsg.GlobalStatusMap()["Created_tmp_tables"])
			if tmpint, _ := strconv.Atoi(v); tmpDiskfileUsageRate > tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "tmpDiskfileUsageRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "tmpDiskfileUsageRate", strconv.Itoa(tmpDiskfileUsageRate))
				result.R.Performance.PerformanceStatus.TmpDiskfileUsageRate = append(result.R.Performance.PerformanceStatus.TmpDiskfileUsageRate, d)
				glog.Error("磁盘临时文件使用过多。检查慢速SQL日志或参数“tmp_table_size”和“max_heap_table_size”")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "tmpDiskfileUsageRate"
				d["checkStatus"] = "normal" // 正常
				d["currentValue"] = fmt.Sprintf("%s=%s", "tmpDiskfileUsageRate", strconv.Itoa(tmpDiskfileUsageRate))
				result.R.Performance.PerformanceStatus.TmpDiskfileUsageRate = append(result.R.Performance.PerformanceStatus.TmpDiskfileUsageRate, d)
			}
		}
		// 统计数据库表扫描率  handler_read_rnd_next/com_select *100
		//	tableScanUsageRate,err := stream.Strea.Percentage(dbmsg.GlobalStatusMap() ["handler_read_rnd_next"],dbmsg.GlobalStatusMap() ["com_select"])
		//	if tableScanUsageRate > 10 && err == nil{
		//		glog.Error("Too many disk temporary file are being used. Check the slow SQL log or parameters tmp_table_size and max_heap_table_size")
		//	}

		// 统计数据库Innodb buffer pool 使用率 100 - (Innodb_buffer_pool_pages_free * 100 / Innodb_buffer_pool_pages_total) # 单位为%
		if strings.EqualFold(k, "innodbBufferPoolUsageRate") {
			innodbBufferPoolUsageRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["Innodb_buffer_pool_pages_free"], dbmsg.GlobalStatusMap()["Innodb_buffer_pool_pages_total"])
			if tmpint, _ := strconv.Atoi(v); 100-innodbBufferPoolUsageRate > tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf("<%s%%", v)
				d["checkType"] = "innodbBufferPoolUsageRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "innodbBufferPoolUsageRate", strconv.Itoa(innodbBufferPoolUsageRate))
				result.R.Performance.PerformanceStatus.InnodbBufferPoolUsageRate = append(result.R.Performance.PerformanceStatus.InnodbBufferPoolUsageRate, d)
				glog.Warn(" PF1-05 The InnoDB buffer pool usage is lower than 80%")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf("<%s%%", v)
				d["checkType"] = "innodbBufferPoolUsageRate"
				d["checkStatus"] = "normal" // 正常
				d["currentValue"] = fmt.Sprintf("%s=%s", "innodbBufferPoolUsageRate", strconv.Itoa(innodbBufferPoolUsageRate))
				result.R.Performance.PerformanceStatus.InnodbBufferPoolUsageRate = append(result.R.Performance.PerformanceStatus.InnodbBufferPoolUsageRate, d)
			}
		}
		// 统计数据库Innodb buffer pool 的脏页率Innodb_buffer_pool_pages_dirty * 100 / Innodb_buffer_pool_pages_total
		if strings.EqualFold(k, "innodbBufferPoolDirtyPagesRate") {
			innodbBufferPoolDirtyPagesRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["Innodb_buffer_pool_pages_dirty"], dbmsg.GlobalStatusMap()["Innodb_buffer_pool_pages_total"])
			if tmpint, _ := strconv.Atoi(v); innodbBufferPoolDirtyPagesRate > tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "innodbBufferPoolDirtyPagesRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "innodbBufferPoolDirtyPagesRate", strconv.Itoa(innodbBufferPoolDirtyPagesRate))
				result.R.Performance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate = append(result.R.Performance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate, d)
				glog.Warn(" PF1-06 The proportion of dirty pages in the MySQL InnoDB buffer pool exceeds 50%")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)

				d["checkType"] = "innodbBufferPoolDirtyPagesRate"
				d["checkStatus"] = "normal" // 正常
				d["currentValue"] = fmt.Sprintf("%s=%s", "innodbBufferPoolDirtyPagesRate", strconv.Itoa(innodbBufferPoolDirtyPagesRate))
				result.R.Performance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate = append(result.R.Performance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate, d)
			}
		}

		// 统计数据库Innodb buffer pool的命中率Innodb_buffer_pool_reads *100 /Innodb_buffer_pool_read_requests
		if strings.EqualFold(k, "innodbBufferPoolHitRate") {
			innodbBufferPoolHitRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["Innodb_buffer_pool_reads"], dbmsg.GlobalStatusMap()["Innodb_buffer_pool_read_requests"])
			if tmpint, _ := strconv.Atoi(v); 100-innodbBufferPoolHitRate < tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf("<%s%%", v)
				d["checkType"] = "innodbBufferPoolHitRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "innodbBufferPoolHitRate", strconv.Itoa(100-innodbBufferPoolHitRate))
				result.R.Performance.PerformanceStatus.InnodbBufferPoolHitRate = append(result.R.Performance.PerformanceStatus.InnodbBufferPoolHitRate, d)
				glog.Warn("MySQL InnoDB缓冲池cache命中率过低。建议增加“innoDB buffer pool”的大小")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf("<%s%%", v)
				d["checkType"] = "innodbBufferPoolHitRate"
				d["checkStatus"] = "normal" // 正常
				d["currentValue"] = fmt.Sprintf("%s=%s", "innodbBufferPoolHitRate", strconv.Itoa(100-innodbBufferPoolHitRate))
				result.R.Performance.PerformanceStatus.InnodbBufferPoolHitRate = append(result.R.Performance.PerformanceStatus.InnodbBufferPoolHitRate, d)
			}
		}

		// 统计数据库文件句柄使用率open_files / open_files_limit * 100% <= 75％
		if strings.EqualFold(k, "openFileUsageRate") {
			openFileUsageRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["open_files"], dbmsg.GlobalVariablesMap()["open_files_limit"])
			if tmpint, _ := strconv.Atoi(v); openFileUsageRate > tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "openFileUsageRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "openFileUsageRate", strconv.Itoa(openFileUsageRate))
				result.R.Performance.PerformanceStatus.OpenFileUsageRate = append(result.R.Performance.PerformanceStatus.OpenFileUsageRate, d)
				glog.Warn(" PF1-08 If the database file handle usage reaches 75%, you are advised to adjust the open_files_LIMIT parameter")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "openFileUsageRate"
				d["checkStatus"] = "normal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "openFileUsageRate", strconv.Itoa(openFileUsageRate))
				result.R.Performance.PerformanceStatus.OpenFileUsageRate = append(result.R.Performance.PerformanceStatus.OpenFileUsageRate, d)
			}
		}

		// 统计数据库表打开缓存率Open_tables *100/table_open_cache
		if strings.EqualFold(k, "openTableCacheUsageRate") {
			openTableCacheUsageRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["open_files"], dbmsg.GlobalVariablesMap()["open_files_limit"])
			if tmpint, _ := strconv.Atoi(v); openTableCacheUsageRate > tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "openTableCacheUsageRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "openTableCacheUsageRate", strconv.Itoa(openTableCacheUsageRate))
				result.R.Performance.PerformanceStatus.OpenTableCacheUsageRate = append(result.R.Performance.PerformanceStatus.OpenTableCacheUsageRate, d)
				glog.Warn(" PF1-09 Database open table cache usage exceeds 80%, you are advised to adjust the table_open_cache parameter")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "openTableCacheUsageRate"
				d["checkStatus"] = "normal" // 正常
				d["currentValue"] = fmt.Sprintf("%s=%s", "openTableCacheUsageRate", strconv.Itoa(openTableCacheUsageRate))
				result.R.Performance.PerformanceStatus.OpenTableCacheUsageRate = append(result.R.Performance.PerformanceStatus.OpenTableCacheUsageRate, d)
			}
		}

		// 统计数据库表缓存溢出使用率Table_open_cache_overflows *100 /(Table_open_cache_hits+Table_open_cache_misses)
		if strings.EqualFold(k, "openTableCacheOverflowsUsageRate") {
			openTableTotal, err := stream.Strea.Add(dbmsg.GlobalStatusMap()["Table_open_cache_hits"], dbmsg.GlobalStatusMap()["Table_open_cache_misses"])
			openTableCacheOverflowsUsageRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["Table_open_cache_overflows"], openTableTotal)
			if tmpint, _ := strconv.Atoi(v); openTableCacheOverflowsUsageRate > tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "openTableCacheOverflowsUsageRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "openTableCacheOverflowsUsageRate", strconv.Itoa(openTableCacheOverflowsUsageRate))
				result.R.Performance.PerformanceStatus.OpenTableCacheOverflowsUsageRate = append(result.R.Performance.PerformanceStatus.OpenTableCacheOverflowsUsageRate, d)
				glog.Warn(" PF1-10 If the tablespace cache overflow usage is greater than 10%, you are advised to adjust the table_open_cache parameter")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "openTableCacheOverflowsUsageRate"
				d["checkStatus"] = "normal" // 正常
				d["currentValue"] = fmt.Sprintf("%s=%s", "openTableCacheOverflowsUsageRate", strconv.Itoa(openTableCacheOverflowsUsageRate))
				result.R.Performance.PerformanceStatus.OpenTableCacheOverflowsUsageRate = append(result.R.Performance.PerformanceStatus.OpenTableCacheOverflowsUsageRate, d)
			}
		}

		// 统计数据库全表扫描的占比率Select_scan *100 /Queries
		if strings.EqualFold(k, "selectScanUsageRate") {
			selectScanUsageRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["Select_scan"], dbmsg.GlobalStatusMap()["Queries"])
			if tmpInt, _ := strconv.Atoi(v); selectScanUsageRate > tmpInt && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "selectScanUsageRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "selectScanUsageRate", strconv.Itoa(selectScanUsageRate))
				result.R.Performance.PerformanceStatus.SelectScanUsageRate = append(result.R.Performance.PerformanceStatus.SelectScanUsageRate, d)
				glog.Warn("数据库不使用索引。如果全表扫描使用率超过“10%”，建议检查SQL慢速")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "selectScanUsageRate"
				d["checkStatus"] = "normal" // 正常
				d["currentValue"] = fmt.Sprintf("%s=%s", "selectScanUsageRate", strconv.Itoa(selectScanUsageRate))
				result.R.Performance.PerformanceStatus.SelectScanUsageRate = append(result.R.Performance.PerformanceStatus.SelectScanUsageRate, d)
			}
		}
		// 统计数据库join语句发生全表扫描占比率Select_full_join *100 /Queries
		if strings.EqualFold(k, "selectfullJoinScanUsageRate") {
			selectfullJoinScanUsageRate, err := stream.Strea.Percentage(dbmsg.GlobalStatusMap()["Select_full_join"], dbmsg.GlobalStatusMap()["Queries"])
			if tmpint, _ := strconv.Atoi(v); selectfullJoinScanUsageRate > tmpint && err == nil {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "selectfullJoinScanUsageRate"
				d["checkStatus"] = "abnormal" // 异常
				d["currentValue"] = fmt.Sprintf("%s=%s", "selectfullJoinScanUsageRate", strconv.Itoa(selectfullJoinScanUsageRate))
				result.R.Performance.PerformanceStatus.SelectfullJoinScanUsageRate = append(result.R.Performance.PerformanceStatus.SelectfullJoinScanUsageRate, d)
				glog.Warn("数据库使用JOIN语句，非驱动表不使用索引。全表扫描使用率大于“10%”。建议检查慢速SQL")
			} else {
				var d = make(map[string]string)
				d["threshold"] = fmt.Sprintf(">%s%%", v)
				d["checkType"] = "selectfullJoinScanUsageRate"
				d["checkStatus"] = "normal" // 正常
				d["currentValue"] = fmt.Sprintf("%s=%s", "selectfullJoinScanUsageRate", strconv.Itoa(selectfullJoinScanUsageRate))
				result.R.Performance.PerformanceStatus.SelectfullJoinScanUsageRate = append(result.R.Performance.PerformanceStatus.SelectfullJoinScanUsageRate, d)
			}
		}
	}
}
