package dbmsg

import (
	"dbcheck/pkg/db"
	"fmt"
	"sync"
)

type informationSchemaKeyColumnUsage struct {
	CONSTRAINT_CATALOG            string  `db:"CONSTRAINT_CATALOG"`
	CONSTRAINT_SCHEMA             string  `db:"CONSTRAINT_SCHEMA"`
	CONSTRAINT_NAME               string  `db:"CONSTRAINT_NAME"`
	TABLE_CATALOG                 string  `db:"TABLE_CATALOG"`
	TABLE_SCHEMA                  string  `db:"TABLE_SCHEMA"`
	TABLE_NAME                    string  `db:"TABLE_NAME"`
	COLUMN_NAME                   string  `db:"COLUMN_NAME"`
	ORDINAL_POSITION              string  `db:"ORDINAL_POSITION"`
	POSITION_IN_UNIQUE_CONSTRAINT string  `db:"POSITION_IN_UNIQUE_CONSTRAINT"`
	REFERENCED_TABLE_SCHEMA       string  `db:"REFERENCED_TABLE_SCHEMA"`
	REFERENCED_TABLE_NAME         *string `db:"REFERENCED_TABLE_NAME"`
	REFERENCED_COLUMN_NAME        *string `db:"REFERENCED_COLUMN_NAME"`
}

var (
	keyColumnUsage     []informationSchemaKeyColumnUsage // 具有约束的键列  INFORMATION_SCHEMA.KEY_COLUMN_USAGE
	keyColumnUsageOnce sync.Once
)

func InformationSchemaKeyColumnUsage() []informationSchemaKeyColumnUsage {
	keyColumnUsageOnce.Do(func() {
		strSql := fmt.Sprintf("select "+
			"CONSTRAINT_SCHEMA, TABLE_NAME, COLUMN_NAME, CONSTRAINT_NAME, "+
			"REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME "+
			"from INFORMATION_SCHEMA.KEY_COLUMN_USAGE "+
			"where CONSTRAINT_SCHEMA not in (%s)", ignoreTableSchema)
		// strSql := fmt.Sprintf("select "+
		// 	"CONSTRAINT_SCHEMA databaseName,TABLE_NAME tableName,COLUMN_NAME columnName,CONSTRAINT_NAME, "+
		// 	"REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME "+
		// 	"from INFORMATION_SCHEMA.KEY_COLUMN_USAGE "+
		// 	"where CONSTRAINT_SCHEMA not in (%s)", ignoreTableSchema)
		db.Select(&keyColumnUsage, strSql)
	})

	return keyColumnUsage
}
