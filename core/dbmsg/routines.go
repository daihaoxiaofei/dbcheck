package dbmsg

import (
	"dbcheck/pkg/db"
	"fmt"
	"sync"
)

type informationSchemaRoutines struct {
	SPECIFIC_NAME            string `db:"SPECIFIC_NAME"`
	ROUTINE_CATALOG          string `db:"ROUTINE_CATALOG"`
	ROUTINE_SCHEMA           string `db:"ROUTINE_SCHEMA"`
	ROUTINE_NAME             string `db:"ROUTINE_NAME"`
	ROUTINE_TYPE             string `db:"ROUTINE_TYPE"`
	DATA_TYPE                string `db:"DATA_TYPE"`
	CHARACTER_MAXIMUM_LENGTH string `db:"CHARACTER_MAXIMUM_LENGTH"`
	CHARACTER_OCTET_LENGTH   string `db:"CHARACTER_OCTET_LENGTH"`
	NUMERIC_PRECISION        string `db:"NUMERIC_PRECISION"`
	NUMERIC_SCALE            string `db:"NUMERIC_SCALE"`
	DATETIME_PRECISION       string `db:"DATETIME_PRECISION"`
	CHARACTER_SET_NAME       string `db:"CHARACTER_SET_NAME"`
	COLLATION_NAME           string `db:"COLLATION_NAME"`
	DTD_IDENTIFIER           string `db:"DTD_IDENTIFIER"`
	ROUTINE_BODY             string `db:"ROUTINE_BODY"`
	ROUTINE_DEFINITION       string `db:"ROUTINE_DEFINITION"`
	EXTERNAL_NAME            string `db:"EXTERNAL_NAME"`
	EXTERNAL_LANGUAGE        string `db:"EXTERNAL_LANGUAGE"`
	PARAMETER_STYLE          string `db:"PARAMETER_STYLE"`
	IS_DETERMINISTIC         string `db:"IS_DETERMINISTIC"`
	SQL_DATA_ACCESS          string `db:"SQL_DATA_ACCESS"`
	SQL_PATH                 string `db:"SQL_PATH"`
	SECURITY_TYPE            string `db:"SECURITY_TYPE"`
	CREATED                  string `db:"CREATED"`
	LAST_ALTERED             string `db:"LAST_ALTERED"`
	SQL_MODE                 string `db:"SQL_MODE"`
	ROUTINE_COMMENT          string `db:"ROUTINE_COMMENT"`
	DEFINER                  string `db:"DEFINER"`
	CHARACTER_SET_CLIENT     string `db:"CHARACTER_SET_CLIENT"`
	COLLATION_CONNECTION     string `db:"COLLATION_CONNECTION"`
	DATABASE_COLLATION       string `db:"DATABASE_COLLATION"`
}

var (
	routines     []informationSchemaRoutines // 存储过程和函数信息  INFORMATION_SCHEMA.routines
	routinesOnce sync.Once
)

func InformationSchemaRoutines() []informationSchemaRoutines {
	routinesOnce.Do(func() {
		strSql := fmt.Sprintf("select ROUTINE_SCHEMA,ROUTINE_NAME,ROUTINE_TYPE,DEFINER,CREATED "+
			"from information_schema.routines where ROUTINE_SCHEMA not in(%s)", ignoreTableSchema)
		db.Select(&routines, strSql)
	})
	return routines
}
