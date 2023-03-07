package dbmsg

import (
	"dbcheck/pkg/db"
	"fmt"
	"sync"
)

// 表提供有关表索引的信息
type informationSchemaStatistics struct {
	TABLE_CATALOG string  `db:"TABLE_CATALOG"`
	TABLE_SCHEMA  string  `db:"TABLE_SCHEMA"`
	TABLE_NAME    string  `db:"TABLE_NAME"`
	NON_UNIQUE    int     `db:"NON_UNIQUE"`
	INDEX_SCHEMA  string  `db:"INDEX_SCHEMA"`
	INDEX_NAME    string  `db:"INDEX_NAME"`
	SEQ_IN_INDEX  string  `db:"SEQ_IN_INDEX"`
	COLUMN_NAME   string  `db:"COLUMN_NAME"`
	COLLATION     string  `db:"COLLATION"`
	CARDINALITY   int     `db:"CARDINALITY"`
	SUB_PART      *int    `db:"SUB_PART"`
	PACKED        *string `db:"PACKED"`
	NULLABLE      string  `db:"NULLABLE"`
	INDEX_TYPE    string  `db:"INDEX_TYPE"`
	COMMENT       string  `db:"COMMENT"`
	INDEX_COMMENT string  `db:"INDEX_COMMENT"`
	IS_VISIBLE    string  `db:"IS_VISIBLE"`
	EXPRESSION    *string `db:"EXPRESSION"`
}

var (
	statistics     []informationSchemaStatistics // 索引的信息  information_schema.STATISTICS
	statisticsOnce sync.Once
)

func InformationSchemaStatistics() []informationSchemaStatistics {
	statisticsOnce.Do(func() {
		strSql := fmt.Sprintf("select * from information_schema.STATISTICS where table_schema not in (%s)", ignoreTableSchema)
		db.Select(&statistics, strSql)
	})
	return statistics
}
