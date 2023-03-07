package dbmsg

import (
	"dbcheck/pkg/db"
	"fmt"
	"sync"
	"time"
)

type informationSchemaTables struct {
	// Table_catalog   string    `db:"Table_catalog"`   // 数据表登记目录
	// Table_schema    string    `db:"Table_schema"`    // 数据表所属的数据库名
	// Table_name      string    `db:"Table_name"`      // 表名称
	// Table_type      string    `db:"Table_type"`      // 表类型[system view|base table]
	// Engine          string    `db:"Engine"`          // 使用的数据库引擎[MyISAM|CSV|InnoDB]
	// Version         string    `db:"Version"`         // 版本，默认值10
	// Row_format      string    `db:"Row_format"`      // 行格式[Compact|Dynamic|Fixed]
	// Table_rows      string    `db:"Table_rows"`      // 表里所存多少行数据
	// Avg_row_length  string    `db:"Avg_row_length"`  // 平均行长度
	// Data_length     int       `db:"Data_length"`     // 数据长度
	// Max_data_length string    `db:"Max_data_length"` // 最大数据长度
	// Index_length    int       `db:"Index_length"`    // 索引长度
	// Data_free       int       `db:"Data_free"`       // 空间碎片
	// Auto_increment  string    `db:"Auto_increment"`  // 做自增主键的自动增量当前值
	// Create_time     time.Time `db:"Create_time"`     // 表的创建时间
	// Update_time     time.Time `db:"Update_time"`     // 表的更新时间
	// Check_time      string    `db:"Check_time"`      // 表的检查时间
	// Table_collation string    `db:"Table_collation"` // 表的字符校验编码集
	// Checksum        string    `db:"Checksum"`        // 校验和
	// Create_options  string    `db:"Create_options"`  // 创建选项
	// Table_comment   string    `db:"Table_comment"`   // 表的注释、备注

	TABLE_CATALOG   string    `db:"TABLE_CATALOG"`   // 数据表登记目录
	TABLE_SCHEMA    string    `db:"TABLE_SCHEMA"`    // 数据表所属的数据库名
	TABLE_NAME      string    `db:"TABLE_NAME"`      // 表名称
	TABLE_TYPE      string    `db:"TABLE_TYPE"`      // 表类型[system view|base table]
	ENGINE          string    `db:"ENGINE"`          // 使用的数据库引擎[MyISAM|CSV|InnoDB]
	VERSION         string    `db:"VERSION"`         // 版本，默认值10
	ROW_FORMAT      string    `db:"ROW_FORMAT"`      // 行格式[Compact|Dynamic|Fixed]
	TABLE_ROWS      string    `db:"TABLE_ROWS"`      // 表里所存多少行数据
	AVG_ROW_LENGTH  string    `db:"AVG_ROW_LENGTH"`  // 平均行长度
	DATA_LENGTH     int       `db:"DATA_LENGTH"`     // 数据长度
	MAX_DATA_LENGTH string    `db:"MAX_DATA_LENGTH"` // 最大数据长度
	INDEX_LENGTH    int       `db:"INDEX_LENGTH"`    // 索引长度
	DATA_FREE       int       `db:"DATA_FREE"`       // 空间碎片
	AUTO_INCREMENT  *string   `db:"AUTO_INCREMENT"`  // 做自增主键的自动增量当前值
	CREATE_TIME     time.Time `db:"CREATE_TIME"`     // 表的创建时间
	UPDATE_TIME     *string   `db:"UPDATE_TIME"`     // 表的更新时间
	CHECK_TIME      *string   `db:"CHECK_TIME"`      // 表的检查时间
	TABLE_COLLATION string    `db:"TABLE_COLLATION"` // 表的字符校验编码集
	CHECKSUM        *string   `db:"CHECKSUM"`        // 校验和
	CREATE_OPTIONS  string    `db:"CREATE_OPTIONS"`  // 创建选项
	TABLE_COMMENT   string    `db:"TABLE_COMMENT"`   // 表的注释、备注
}

// def	mysql	columns_priv	BASE TABLE	InnoDB	10	Dynamic	0	0	16384	0	0	10485760		2022-03-07 18:32:38			utf8_bin		row_format=DYNAMIC stats_persistent=0	Column privileges
// def	mysql	component	BASE TABLE	InnoDB	10	Dynamic	0	0	16384	0	0	10485760	1	2022-03-07 18:32:38			utf8_general_ci		row_format=DYNAMIC	Components
// def	mysql	db	BASE TABLE	InnoDB	10	Dynamic	2	8192	16384	0	16384	10485760		2022-03-07 18:32:38			utf8_bin		row_format=DYNAMIC stats_persistent=0	Database privileges
// def	mysql	default_roles	BASE TABLE	InnoDB	10	Dynamic	0	0	16384	0	0	10485760		2022-03-07 18:32:38			utf8_bin		row_format=DYNAMIC stats_persistent=0	Default roles
// def	mysql	engine_cost	BASE TABLE	InnoDB	10	Dynamic	2	8192	16384	0	0	10485760		2022-03-07 18:32:38			utf8_general_ci		row_format=DYNAMIC stats_persistent=0

var (
	tables     []informationSchemaTables // 表   information_schema.tables
	tablesOnce sync.Once
)

func InformationSchemaTables() []informationSchemaTables {
	tablesOnce.Do(func() {
		strSql := fmt.Sprintf("select * from information_schema.TABLES where table_schema not in (%s)", ignoreTableSchema)
		db.Select(&tables, strSql)
	})
	return tables
}
