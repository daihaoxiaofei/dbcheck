package inspection

import "C"
import (
	"dbcheck/pkg/db"
	"fmt"
	"go.uber.org/zap"
	"strings"

	"dbcheck/core/dbmsg"
	"dbcheck/core/result"
	"dbcheck/pkg/config"
	"dbcheck/pkg/glog"
)

type GlobalVariablesMap struct {
	Name  string
	Value string
}

// DatabaseConfigCheck 配置参数检查功能
func DatabaseConfigCheck() {
	ConfigurationCanCheck := map[string]string{
		`super_read_only`:                `off`,
		`innodb_read_only`:               `off`,
		`binlog_format`:                  `row`,
		`character_set_server`:           `utf8mb4`,
		`default_authentication_plugin`:  `mysql_native_password`,
		`default_storage_engine`:         `innodb`,
		`default_tmp_storage_engine`:     `innodb`,
		`innodb_flush_log_at_trx_commit`: `1`,
		`innodb_flush_method`:            `fsync`,
		`innodb_deadlock_detect`:         `on`,
		`relay_log_purge`:                `on`,
		`sync_binlog`:                    `1`,
		`system_time_zone`:               `CST`,
		`time_zone`:                      `system`,
		// 事务隔离级别 读未提交(READ-UNCOMMITTED) 读已提交(READ-COMMITTED) 可重复读(REPEATABLE-READ) 序列化(SERIALIZABLE)
		// 级别越高， 多个事务在并发访问数据库时互相产生数据干扰的可能性越低，但是并发访问的性能就越差。
		`transaction_isolation`: `REPEATABLE-READ`,
		`transaction_read_only`: `off`,
		`unique_checks`:         `on`,
		// `read_only`:             `true`, // 从库需要
		// `relay_log_recovery`:    `on`,   // 主从复制时 保证宕机时数据可靠性
	}

	for k, v := range ConfigurationCanCheck {
		a, ok := dbmsg.GlobalVariablesMap()[k]
		if !ok {
			glog.Error("当前数据配置参数不存在。请检查是否打错了", zap.String(`i`, k))
			continue
		}
		d := make(map[string]string)
		d["configVariableName"] = k
		d["configVariable"] = a     // 当前值
		d["configValue"] = v        // 建议值
		d["checkStatus"] = "normal" // 正常
		d["checkType"] = "configParameter"
		if !strings.EqualFold(a, v) {
			d["checkStatus"] = "abnormal" // 异常
			d["threshold"] = v
			d["currentValue"] = fmt.Sprintf("%s=%s", k, a)
			errorStrInfo := fmt.Sprintf("GlobalVariables 配置项: %s 不符合预定要求 当前值: %s 建议设置为: %s", k, a, v)
			glog.Error(errorStrInfo)
		}
		result.R.Config.ConfigParameter = append(result.R.Config.ConfigParameter, d)
	}
}

type DatabaseBaselineCheckStruct struct{}

type TableDesignComplianceStruct struct {
	CONSTRAINT_SCHEMA interface{} `json:"CONSTRAINT_SCHEMA"`
	TableName         interface{} `json:"tableName"`
	Engine            interface{} `json:"engine"`
	Charset           interface{} `json:"charset"`
}

func newMap(source map[string]string) map[string]string {
	var n = make(map[string]string)
	for k, v := range source {
		n[k] = v
	}
	return n
}
func newD(schema, table string) map[string]string {
	return map[string]string{
		"database":    schema,
		"tableName":   table,
		"checkStatus": "normal",
	}
}

var BaselineCanCheck = map[string]string{
	`tableAutoIncrement`:        `bigint`,                     // 检测自增主键数据类型 		自增主键非bigint类型
	`tableBigColumns`:           `text,timestamp,blob`,        // 检测表是否存在大列类型 	表中是否存在大字段类型（text,timestamp,blob）
	`tableIncludeRepeatIndex`:   ``,                           // 检测表中是否存在冗余索引
	`tableProcedureFuncTrigger`: `procedure,function,trigger`, // 检测存储过程、函数、触发器使用情况	判断是否使用存储过程、存储函数、触发器
}

// tableNoPrimaryKey
// 没有主键的表 按顺序排列的 如果第一个字段不是主键则后面就没有主键了
// 主键自增列是否为 bigint
// 表中是否存在大字段 blob、text、varchar(8099)、timestamp 数据类型 非自增时判断
func tableNoPrimaryKey() {
	var (
		lastSchema, lastTable string
		columnNum             int
	)
	for _, v := range dbmsg.InformationSchemaColumns() {
		d := newD(v.TABLE_SCHEMA, v.TABLE_NAME)
		d["checkType"] = "tableNoPrimaryKey"
		if lastSchema != v.TABLE_SCHEMA || lastTable != v.TABLE_NAME { // 是否进入新表
			if v.COLUMN_KEY != "PRI" { // 不是主键 &&
				d["checkStatus"] = "abnormal" // 异常
				d["threshold"] = fmt.Sprintf("%s", ``)
				d["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
				glog.Error(fmt.Sprintf("当前表没有主键, Schema: %s table:  %s", v.TABLE_SCHEMA, v.TABLE_NAME))
			}
			if columnNum > 255 {
				glog.Error(fmt.Sprintf("当前表列数大于255, Schema: %s table:  %s", v.TABLE_SCHEMA, v.TABLE_NAME))
			}
			columnNum = 0
		} else {
			columnNum++
		}

		lastSchema = v.TABLE_SCHEMA
		lastTable = v.TABLE_NAME
		result.R.Baseline.TableDesign.TableNoPrimaryKey = append(result.R.Baseline.TableDesign.TableNoPrimaryKey, d)

		c := newD(v.TABLE_SCHEMA, v.TABLE_NAME)
		c["columnName"] = v.COLUMN_NAME
		// 主键自增列是否为 bigint
		if v.EXTRA == "auto_increment" {
			c["checkType"] = "tableAutoIncrement"
			if !strings.Contains(v.COLUMN_TYPE, `bigint`) {
				c["checkStatus"] = "abnormal" // 异常
				c["threshold"] = "无自增主键"
				c["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
				glog.Error(fmt.Sprintf("主键列不是Bigint类型 database: %s tableName: %s columnsName: %s columnType: %s.",
					v.TABLE_SCHEMA, v.TABLE_NAME, v.COLUMN_NAME, v.COLUMN_TYPE))
				if config.C.Repair {
					db.Exec(fmt.Sprintf("alter table %s.%s modify column %s BIGINT ", v.TABLE_SCHEMA, v.TABLE_NAME, v.COLUMN_NAME))
					glog.Error(fmt.Sprintf("已修复主键为Bigint类型 database: %s tableName: %s columnName: %s",
						v.TABLE_SCHEMA, v.TABLE_NAME, v.COLUMN_NAME))
				}
			}
			result.R.Baseline.ColumnDesign.TableAutoIncrement = append(result.R.Baseline.ColumnDesign.TableAutoIncrement, d)
		} else {
			// 表中是否存在大字段 blob、text、varchar(8099)、timestamp 数据类型 非自增时判断
			m := newD(v.TABLE_SCHEMA, v.TABLE_NAME)
			m["checkType"] = "tableBigColumns"
			if strings.Contains(`text,timestamp,blob`, v.DATA_TYPE) {
				m["checkStatus"] = "abnormal" // 异常
				m["threshold"] = fmt.Sprintf("%s", `text,timestamp,blob`)
				m["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
				glog.Error(fmt.Sprintf("数据库中有不合理类型: %s database: %s tableName: %s columnsName: %s",
					v.COLUMN_TYPE, v.TABLE_SCHEMA, v.TABLE_NAME, v.COLUMN_NAME))
			}
			result.R.Baseline.ColumnDesign.TableBigColumns = append(result.R.Baseline.ColumnDesign.TableBigColumns, m)
		}
	}

}

// BaselineCheckIndexColumnDesign 索引设计合规性
// 索引列可为空情况
// 索引列是enum、set、blob或文本类型
// 寻找冗余索引
func BaselineCheckIndexColumnDesign() {
	var indexesMap = make(map[string]string)
	for _, v := range dbmsg.InformationSchemaStatistics() {
		vName := fmt.Sprintf("%s_%s_%s", v.TABLE_SCHEMA, v.TABLE_NAME, v.COLUMN_NAME)
		indexesMap[vName] = v.INDEX_NAME
	}
	for _, v := range dbmsg.InformationSchemaColumns() {
		d := newD(v.TABLE_SCHEMA, v.TABLE_NAME)
		d["columnName"] = v.COLUMN_NAME
		d["columnType"] = v.COLUMN_TYPE
		d["columnIsNull"] = v.IS_NULLABLE
		vName := fmt.Sprintf("%s_%s_%s", v.TABLE_SCHEMA, v.TABLE_NAME, v.COLUMN_NAME)
		// 该字段存在索引
		if _, ok := indexesMap[vName]; ok {
			d["checkStatus"] = "normal" // 异常
			d["checkType"] = "indexColumnIsNull"
			// 判断索引列是否允许为空
			if strings.EqualFold(`YES`, v.IS_NULLABLE) {
				d["checkStatus"] = "abnormal" // 异常
				d["threshold"] = "索引列可为空"
				d["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
				glog.Error(fmt.Sprintf("索引列可为空 database: %s  tablename: %s indexName: %s columnName: %s columnType: %s",
					v.TABLE_SCHEMA, v.TABLE_NAME, indexesMap[vName], v.COLUMN_NAME, v.COLUMN_TYPE))
				if config.C.Repair {
					db.Exec(fmt.Sprintf("ALTER TABLE %s.%s CHANGE COLUMN %s %s %s NOT NULL",
						v.TABLE_SCHEMA, v.TABLE_NAME, v.COLUMN_NAME, v.COLUMN_NAME, v.COLUMN_TYPE))
					glog.Error(fmt.Sprintf("尝试修正 索引列可为空情况 database:database: %s  tablename: %s indexName: %s columnName: %s columnType: %s",
						v.TABLE_SCHEMA, v.TABLE_NAME, indexesMap[vName], v.COLUMN_NAME, v.COLUMN_TYPE))
				}
			}
			result.R.Baseline.IndexColumnsDesign.IndexColumnIsNull = append(result.R.Baseline.IndexColumnsDesign.IndexColumnIsNull, d)

			// 判断索引列是否建立在enum、set、blob、text类型上面
			tmpColumnType := strings.ToLower(v.DATA_TYPE)
			m := newMap(d)
			m["checkStatus"] = "normal"
			m["checkType"] = "indexColumnType"
			if strings.Contains(`enum,set,blob,text`, tmpColumnType) {
				m["checkStatus"] = "abnormal" // 异常
				m["threshold"] = fmt.Sprintf("%s", `enum,set,blob,text`)
				m["currentValue"] = fmt.Sprintf("%s.%s", v.TABLE_SCHEMA, v.TABLE_NAME)
				glog.Error(fmt.Sprintf("索引列是enum、set、blob或文本类型 database: %s  tablename: %s indexName: %s columnName: %s columnType: %s",
					v.TABLE_SCHEMA, v.TABLE_NAME, indexesMap[vName], v.COLUMN_NAME, v.COLUMN_TYPE))
			}
			result.R.Baseline.IndexColumnsDesign.IndexColumnType = append(result.R.Baseline.IndexColumnsDesign.IndexColumnType, m)
		}
	}

	// 寻找冗余索引
	// 利用map合并联合索引列
	var tmpIndexMargeMap = make(map[string]string)
	b := dbmsg.InformationSchemaStatistics()
	for k, v := range dbmsg.InformationSchemaStatistics() {
		key := fmt.Sprintf("%s@%s@@%s", v.TABLE_SCHEMA, v.TABLE_NAME, v.INDEX_NAME)
		if val, ok := tmpIndexMargeMap[key]; ok && k > 1 && b[k-1].TABLE_SCHEMA == b[k].TABLE_SCHEMA && b[k-1].TABLE_NAME == b[k].TABLE_NAME {
			tmpValue := fmt.Sprintf("%s,%s", val, b[k].COLUMN_NAME)
			tmpIndexMargeMap[key] = tmpValue
		} else {
			tmpValue := fmt.Sprintf("%s", b[k].COLUMN_NAME)
			tmpIndexMargeMap[key] = tmpValue
		}
	}
	// 分离出每个库表下包含的索引
	var tableIncludeIndexMap = make(map[string]map[string]string)
	for k, v := range tmpIndexMargeMap {
		tmpMap := make(map[string]string)
		a := strings.Split(k, "@@") // 库表
		if val, ok := tableIncludeIndexMap[a[0]]; ok {
			for tmpK := range val {
				tmpMap[tmpK] = val[tmpK] // 旧的key value
			}
			tmpMap[a[1]] = v // 新的key value
			tableIncludeIndexMap[a[0]] = tmpMap
		} else {
			tmpMap[a[1]] = v
			tableIncludeIndexMap[a[0]] = tmpMap
		}
	}
	// 遍历每一个库表下的索引列，寻找冗余索引
	for k, v := range tableIncludeIndexMap {
		var isRedundancy bool // 多余
		var tmpDatabase, tmpTableName, tmpIndexRedundancyName, tmpIndexRedundancyColumn, tmpIndexColumnName, tmpIndexIncludeColumn string
		a := strings.Split(k, "@")
		tmpDatabase = a[0]
		tmpTableName = a[1]
		for ki, ui := range v {
			for kii, uii := range v {
				if ui != uii && strings.HasPrefix(uii, ui) {
					tmpIndexRedundancyColumn = ui
					tmpIndexIncludeColumn = uii
					tmpIndexColumnName = kii
					tmpIndexRedundancyName = ki
					isRedundancy = true
				}
			}
		}

		d := newD(tmpDatabase, tmpTableName)
		d["redundantIndexes"] = fmt.Sprintf("%s %s,%s %s", tmpIndexRedundancyName, tmpIndexRedundancyColumn, tmpIndexColumnName, tmpIndexIncludeColumn)
		d["checkType"] = "tableIncludeRepeatIndex"
		if isRedundancy {
			d["checkStatus"] = "abnormal" // 异常
			d["threshold"] = "存在重复索引"
			d["currentValue"] = fmt.Sprintf("%s.%s", tmpDatabase, tmpTableName)
			glog.Error(fmt.Sprintf("出现冗余索引列 database:%s tablename: %s Redundant indexes: (indexName:%s indexColumns %s), (indexName: %s indexColumns: %s)",
				tmpDatabase, tmpTableName, tmpIndexRedundancyName, tmpIndexRedundancyColumn, tmpIndexColumnName, tmpIndexIncludeColumn))
		}
		result.R.Baseline.IndexColumnsDesign.IndexColumnIsRepeatIndex = append(result.R.Baseline.IndexColumnsDesign.IndexColumnIsRepeatIndex, d)
	}

}

// BaselineCheckProcedureTriggerDesign 存储过程、存储函数、触发器检查限制
func BaselineCheckProcedureTriggerDesign() {
	if vi, okk := BaselineCanCheck["tableprocedurefunctrigger"]; okk {
		for _, v := range dbmsg.InformationSchemaRoutines() {
			var d = make(map[string]string)
			d["database"] = v.ROUTINE_SCHEMA
			if strings.Contains(strings.ToLower(vi), strings.ToLower(v.ROUTINE_TYPE)) {
				d["checkStatus"] = "abnormal" // 异常状态
				d["checkType"] = "tableProcedureFunc"
				d["threshold"] = fmt.Sprintf("%s", strings.ToLower(v.ROUTINE_TYPE))
				d["currentValue"] = fmt.Sprintf("%s.%s", v.ROUTINE_SCHEMA, v.ROUTINE_NAME)
				result.R.Baseline.ProcedureTriggerDesign.TableProcedure = append(result.R.Baseline.ProcedureTriggerDesign.TableProcedure, d)
				glog.Error(fmt.Sprintf(" BL4-PT01 The current database uses a storage procedure or storage func. The information is as follows: database: %s routineName: %s user: %s create time: %s",
					v.ROUTINE_SCHEMA, v.ROUTINE_NAME, v.DEFINER, v.CREATED))
			}
		}
		// 检查是否使用触发器
		for _, v := range dbmsg.InformationSchemaTriggers() {
			var d = make(map[string]string)
			d["database"] = v.TRIGGER_SCHEMA
			if strings.Contains(strings.ToLower(vi), "trigger") && v.TRIGGER_NAME != `` {
				d["checkStatus"] = "abnormal" // 异常状态
				d["checkType"] = "tableTrigger"
				d["threshold"] = fmt.Sprintf("%s", "trigger")
				d["currentValue"] = fmt.Sprintf("%s.%s", v.TRIGGER_SCHEMA, v.TRIGGER_NAME)
				result.R.Baseline.ProcedureTriggerDesign.TableTrigger = append(result.R.Baseline.ProcedureTriggerDesign.TableTrigger, d)
				glog.Error(fmt.Sprintf(" BL4-PT02 The current database uses a trigger. The information is as follows: database: %s triggerName: %s  user: %s  create time:%s",
					v.TRIGGER_SCHEMA, v.TRIGGER_NAME, v.DEFINER, v.CREATED))
			}
		}
	}
}
