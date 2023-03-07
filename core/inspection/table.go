package inspection

import (
	"dbcheck/core/dbmsg"
	"dbcheck/core/result"
	"dbcheck/core/stream"
	"dbcheck/pkg/db"
	"dbcheck/pkg/glog"
	"fmt"
	"strconv"
	"strings"
)

// DatabasePerformanceTableIndexCheck 检查 information_schema.TABLES
// 表字符集检查 utf8mb4
// 检查是否存在自增主键溢出风险。统计数据库自增id列快要溢出的表
// 单表行数大于500w，且平均行长大于10KB。
// 单表大于6G，并且碎片率大于30%。
func DatabasePerformanceTableIndexCheck() {
	var nowDateTime string
	db.Get(&nowDateTime, "select now() as datetime")

	// nowDateTime := publicClass.DBExecInter.DBQueryDateString()
	var tmpTableColumnMap = make(map[string]string)
	for _, v := range dbmsg.InformationSchemaColumns() {
		// 过滤自增主键为int类型的库表，并输出库、表、数据类型等信息
		if v.EXTRA == "auto_increment" && v.DATA_TYPE == "int" {
			keyDT := fmt.Sprintf("%s@%s", v.TABLE_SCHEMA, v.TABLE_NAME)
			valCT := fmt.Sprintf("%s@%s@%s", v.COLUMN_NAME, v.DATA_TYPE, v.COLUMN_TYPE)
			tmpTableColumnMap[keyDT] = valCT
		}
	}
	CollationsMap := dbmsg.CollationsMap() // 字符集对应表
	for _, v := range dbmsg.InformationSchemaTables() {
		// 表的字符校验编码集
		tableCharset := CollationsMap[v.TABLE_COLLATION]

		// 表字符集检查 utf8mb4
		var d = make(map[string]string)
		d["database"] = v.TABLE_SCHEMA
		d["tableName"] = v.TABLE_NAME
		d["charset"] = tableCharset
		d["checkType"] = "tableCharset"
		d["checkStatus"] = "normal"
		if !strings.Contains(`utf8mb4`, tableCharset) {
			d["checkStatus"] = "abnormal" // 异常
			d["threshold"] = fmt.Sprintf("非%s", `utf8mb4`)
			d["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
			glog.Error(fmt.Sprintf(" 当前表字符集不是 UTF8MB4  Database: %s table: %s table charset: %s ",
				v.TABLE_SCHEMA, v.TABLE_NAME, tableCharset))
		}
		result.R.Baseline.TableDesign.TableCharset = append(result.R.Baseline.TableDesign.TableCharset, d)

		// 检查引擎 innodb
		m := newMap(d)
		m["checkType"] = "tableEngine"
		if v.ENGINE != `` && !strings.EqualFold(v.ENGINE, "innodb") {
			m["checkStatus"] = "abnormal"
			m["threshold"] = fmt.Sprintf("非%s", `tableEngine`)
			m["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
			glog.Error(fmt.Sprintf("当前表引擎集不是innodb引擎 Database: %s table: %s engine: %s ",
				v.TABLE_SCHEMA, v.TABLE_NAME, v.ENGINE))
		}
		if v.ENGINE != `` && strings.EqualFold(v.ENGINE, "innodb") {
			m["checkStatus"] = "normal"
			m["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
		}
		result.R.Baseline.TableDesign.TableEngine = append(result.R.Baseline.TableDesign.TableEngine, m)

		for ki, vi := range PerformanceCanCheck {
			databaseTableName := fmt.Sprintf("%s@%s", v.TABLE_SCHEMA, v.TABLE_NAME)
			if strings.EqualFold(ki, "tableAutoPrimaryKeyUsageRate") {
				if _, ok := tmpTableColumnMap[databaseTableName]; ok {
					e := newMap(d)
					tmpColumnInfoSliect := strings.Split(tmpTableColumnMap[databaseTableName], "@")
					tmpColumnName := tmpColumnInfoSliect[0]
					tmpColumnType := tmpColumnInfoSliect[2]
					e["columnName"] = tmpColumnName
					e["columnType"] = tmpColumnType
					e["threshold"] = fmt.Sprintf(">%s%%", vi)
					e["checkType"] = "tableAutoPrimaryKeyUsageRate"
					// e["autoIncrement"] = strconv.Itoa(int(v.Auto_increment.(int64))) // xiaofei
					// 检查是否存在自增主键溢出风险。统计数据库自增id列快要溢出的表
					if strings.Contains(tmpColumnType, "unsigned") {
						unsignedIntUsageRate, err := stream.Strea.Percentage(v.AUTO_INCREMENT, MunsignedInt)
						if tmpint, _ := strconv.Atoi(vi); unsignedIntUsageRate > tmpint && err == nil {
							e["checkStatus"] = "abnormal"
							e["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
							glog.Warn(fmt.Sprintf(" PF2-01 The self-value-added usage of tables in the database exceeds 85%%, causing data type overflow risks. The details are as follows: Database: %v, table name: %v, increment column name: %v, increment column data type: %v, current increment values: %v",
								v.TABLE_SCHEMA, v.TABLE_NAME, tmpColumnName, tmpColumnType, v.AUTO_INCREMENT))
						} else {
							e["checkStatus"] = "normal"
						}
						result.R.Performance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate = append(result.R.Performance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate, e)
					} else {
						intUsageRate, err := stream.Strea.Percentage(v.AUTO_INCREMENT, Mint)
						if tmpint, _ := strconv.Atoi(vi); intUsageRate > tmpint && err == nil {
							e["checkStatus"] = "abnormal"
							glog.Warn(fmt.Sprintf(" PF2-01 The self-value-added usage of tables in the database exceeds 85%%, causing data type overflow risks. The details are as follows: Database: %v, table name: %v, increment column name: %v, increment column data type: %v, current increment values: %v",
								v.TABLE_SCHEMA, v.TABLE_NAME, tmpColumnName, tmpColumnType, v.AUTO_INCREMENT))
						} else {
							e["checkStatus"] = "normal"
						}
						result.R.Performance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate = append(result.R.Performance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate, e)
					}
				}
			}
			// 单表行数大于500w，且平均行长大于10KB。
			if strings.EqualFold(ki, "tableRows") {
				m := newMap(d)
				tableRows, _ := strconv.Atoi(fmt.Sprintf("%s", v.TABLE_ROWS))
				avgRowLength, _ := strconv.Atoi(fmt.Sprintf("%s", v.AVG_ROW_LENGTH))
				tmpTableRowsThreshold, _ := strconv.Atoi(strings.Split(vi, ",")[0])
				tmpTableAvgRowLength, _ := strconv.Atoi(strings.Split(vi, ",")[1])
				m["threshold"] = fmt.Sprintf(">%dW", tmpTableRowsThreshold/10000)
				m["checkType"] = "tableRows"
				m["tableRows"] = strconv.Itoa(tableRows)
				m["avgRowLength"] = strconv.Itoa(avgRowLength)
				if tableRows > tmpTableRowsThreshold && avgRowLength/1024 > tmpTableAvgRowLength {
					m["checkStatus"] = "abnormal"
					m["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
					result.R.Performance.PerformanceTableIndex.TableRows = append(result.R.Performance.PerformanceTableIndex.TableRows, m)
					glog.Warn(fmt.Sprintf(" PF2-02 The current table is a large table if the number of rows is greater than 5 million and the average line length is greater than 10KB. The details are as follows: Database: %v, table name: %v, tableRows: %v, avgRowLength:%d", v.TABLE_SCHEMA, v.TABLE_NAME, tableRows, avgRowLength/1024))
				} else {
					m["checkStatus"] = "normal"
					result.R.Performance.PerformanceTableIndex.TableRows = append(result.R.Performance.PerformanceTableIndex.TableRows, m)
				}
			}
			var dataLength, indexLength, dataFree int
			if v.DATA_LENGTH != 0 {
				dataLength = v.DATA_LENGTH
			}
			if v.INDEX_LENGTH != 0 {
				indexLength = v.INDEX_LENGTH
			}
			if v.DATA_FREE != 0 {
				dataFree = v.DATA_FREE
			}
			// 单表大于6G，并且碎片率大于30%。
			if strings.EqualFold(ki, "diskFragmentationRate") {
				n := newMap(d)
				tmpDigestTableSizeThreshold, _ := strconv.Atoi(strings.Split(vi, ",")[0])
				tmpDigestTableDiskFragmentationRateThreshold, _ := strconv.Atoi(strings.Split(vi, ",")[1])
				dataLengthTotal := dataLength + indexLength // 表空间
				n["threshold"] = fmt.Sprintf(">%d%%", tmpDigestTableDiskFragmentationRateThreshold)
				n["checkType"] = "diskFragmentationRate"
				if diskFragmentationRate, err := stream.Strea.Percentage(dataFree, dataLengthTotal); diskFragmentationRate > tmpDigestTableDiskFragmentationRateThreshold && err == nil && dataLengthTotal/1024/1024/1024 > tmpDigestTableSizeThreshold {
					n["checkStatus"] = "abnormal"
					n["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
					result.R.Performance.PerformanceTableIndex.DiskFragmentationRate = append(result.R.Performance.PerformanceTableIndex.DiskFragmentationRate, n)
					glog.Warn(fmt.Sprintf(" PF2-03 If the current tablespace contains more than 6 GB and the disk fragmentation rate is greater than 30%%, you are advised to run THE ALTER command to delete disk fragmentation. The details are as follows: Database: %v, table name: %v, Table space size: %dG, diskFragmentationRate:%d", v.TABLE_SCHEMA, v.TABLE_NAME, dataLengthTotal/1024/1024/1024, diskFragmentationRate))
				} else {
					n["diskFragmentationRate"] = strconv.Itoa(diskFragmentationRate)
					n["checkStatus"] = "normal"
					result.R.Performance.PerformanceTableIndex.DiskFragmentationRate = append(result.R.Performance.PerformanceTableIndex.DiskFragmentationRate, n)
				}
			}

			// 单表行数大于1000W，且表空间大于30G
			if strings.EqualFold(ki, "bigTable") {
				dataLengthTotal := dataLength + indexLength // 表空间
				tableRows, _ := strconv.Atoi(fmt.Sprintf("%s", v.TABLE_ROWS))
				tmpDigestTableRowsThreshold, _ := strconv.Atoi(strings.Split(vi, ",")[0])
				tmpDigestTableSizeThreshold, _ := strconv.Atoi(strings.Split(vi, ",")[1])
				o := newMap(d)
				o["threshold"] = fmt.Sprintf(">%dG", tmpDigestTableSizeThreshold)
				o["checkType"] = "bigTable"
				o["dataLengthTotal"] = strconv.Itoa(dataLengthTotal / 1024 / 1024 / 1024)
				o["tableRows"] = strconv.Itoa(tableRows)
				if dataLengthTotal/1024/1024/1024 > tmpDigestTableSizeThreshold && tableRows > tmpDigestTableRowsThreshold {
					o["checkStatus"] = "abnormal"
					o["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
					result.R.Performance.PerformanceTableIndex.BigTable = append(result.R.Performance.PerformanceTableIndex.BigTable, o)
					glog.Warn(fmt.Sprintf(" PF2-04 If the number of rows in the current table is greater than 1000W and the tablespace is greater than 30G, the table belongs to a large table. Recommended Attention Table. The details are as follows: Database: %v, table name: %v, tableRows:%d, Table space size: %dG", v.TABLE_SCHEMA, v.TABLE_NAME, tableRows, dataLengthTotal/1024/1024/1024))
				} else {
					o["checkStatus"] = "normal"
					result.R.Performance.PerformanceTableIndex.BigTable = append(result.R.Performance.PerformanceTableIndex.BigTable, o)
				}
			}

			// 检查一个星期内未更新的表
			if strings.EqualFold(ki, "coldTable") {
				var tableUpdateTime string
				if v.CREATE_TIME.Second() != 0 {
					if v.CREATE_TIME.Second() != 0 {
						tableUpdateTime = v.CREATE_TIME.Format(`2006-01-02 15:04:05`)
					} else {
						tableUpdateTime = *v.UPDATE_TIME
					}
				}
				arrDay, _ := stream.Strea.GetTimeDayArr(tableUpdateTime, nowDateTime)
				p := newMap(d)
				p["threshold"] = fmt.Sprintf(">%s day", vi)
				p["checkType"] = "coldTable"
				if tmpint, _ := strconv.Atoi(vi); arrDay > int64(tmpint) {
					p["checkStatus"] = "abnormal"
					p["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
					result.R.Performance.PerformanceTableIndex.ColdTable = append(result.R.Performance.PerformanceTableIndex.ColdTable, p)
					glog.Error(fmt.Sprintf(" PF2-05 The current table has not been updated for seven days (no DML has occurred against the table). The details are as follows: Database: %v, table name: %v, lasterUpdateTime:%v ", v.TABLE_SCHEMA, v.TABLE_NAME, tableUpdateTime))
				} else {
					p["checkStatus"] = "normal"
					result.R.Performance.PerformanceTableIndex.ColdTable = append(result.R.Performance.PerformanceTableIndex.ColdTable, p)
				}
			}
		}
	}
}
