package inspection

import (
	"fmt"

	"dbcheck/core/dbmsg"
	"dbcheck/core/result"
	"dbcheck/pkg/glog"
)

// InformationSchemaKeyColumnUsage 检测是否存在外键约束
func InformationSchemaKeyColumnUsage() {
	for _, v := range dbmsg.InformationSchemaKeyColumnUsage() {
		d := newD(v.CONSTRAINT_SCHEMA, v.TABLE_NAME)
		d["checkType"] = "tableForeign"
		if v.REFERENCED_TABLE_NAME != nil && v.REFERENCED_COLUMN_NAME != nil {
			d["checkStatus"] = "abnormal" // 异常
			d["threshold"] = fmt.Sprintf("%s", ``)
			d["currentValue"] = fmt.Sprintf("%s.%s", v.CONSTRAINT_SCHEMA, v.TABLE_NAME)
			d["columnName"] = v.COLUMN_NAME
			d["constraintName"] = v.CONSTRAINT_NAME
			d["referencedTableName"] = *v.REFERENCED_TABLE_NAME
			d["referencedColumnName"] = *v.REFERENCED_COLUMN_NAME
			glog.Error(fmt.Sprintf(" BL1-TC03 The current table uses a foreign key constraint. The information is as follows: database: %s "+
				"tableName: %s column: %s Foreign key constraint name: %s Foreign key constraints table: %s"+
				"Foreign key constraints columns: %s",
				v.CONSTRAINT_SCHEMA, v.TABLE_NAME, v.COLUMN_NAME, v.CONSTRAINT_NAME, *v.REFERENCED_TABLE_NAME, *v.REFERENCED_COLUMN_NAME))
		}
		result.R.Baseline.TableDesign.TableForeign = append(result.R.Baseline.TableDesign.TableForeign, d)
	}
}
