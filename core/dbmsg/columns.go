package dbmsg

import (
	"dbcheck/pkg/db"
	"fmt"
	"sync"
)

type informationSchemaColumns struct {
	// TABLE_CATALOG            string `db:"TABLE_CATALOG"`            // 表限定符
	// TABLE_SCHEMA             string `db:"TABLE_SCHEMA"`             // 表所有者
	// TABLE_NAME               string `db:"TABLE_NAME"`               // 表名
	// COLUMN_NAME              string `db:"COLUMN_NAME"`              // 列名
	// ORDINAL_POSITION         string `db:"ORDINAL_POSITION"`         // 列标识号
	// COLUMN_DEFAULT           string `db:"COLUMN_DEFAULT"`           // 列的默认值
	// IS_NULLABLE              string `db:"IS_NULLABLE"`              // 列的为空性。如果列允许 NULL，那么该列返回 YES。否则，返回 NO
	// DATA_TYPE                string `db:"DATA_TYPE"`                // 系统提供的数据类型
	// CHARACTER_MAXIMUM_LENGTH string `db:"CHARACTER_MAXIMUM_LENGTH"` // 以字符为单位的最大长度，适于二进制数据、字符数据，或者文本和图像数据。否则，返回 NULL。有关更多信息，请参见数据类型
	// CHARACTER_OCTET_LENGTH   string `db:"CHARACTER_OCTET_LENGTH"`   // 以字节为单位的最大长度，适于二进制数据、字符数据，或者文本和图像数据。否则，返回 NULL
	// NUMERIC_PRECISION        string `db:"NUMERIC_PRECISION"`        // 近似数字数据、精确数字数据、整型数据或货币数据的精度。否则，返回 NULL
	// NUMERIC_PRECISION_RADIX  string `db:"NUMERIC_PRECISION_RADIX"`  // 近似数字数据、精确数字数据、整型数据或货币数据的精度基数。否则，返回 NULL
	// NUMERIC_SCALE            string `db:"NUMERIC_SCALE"`            // 近似数字数据、精确数字数据、整数数据或货币数据的小数位数。否则，返回 NULL
	// DATETIME_PRECISION       string `db:"DATETIME_PRECISION"`       // datetime 及 SQL-92 interval 数据类型的子类型代码。对于其它数据类型，返回 NULL
	// CHARACTER_SET_CATALOG    string `db:"CHARACTER_SET_CATALOG"`    // 如果列是字符数据或 text 数据类型，那么返回 master，指明字符集所在的数据库。否则，返回 NULL
	// CHARACTER_SET_SCHEMA     string `db:"CHARACTER_SET_SCHEMA"`     // 如果列是字符数据或 text 数据类型，那么返回 DBO，指明字符集的所有者名称。否则，返回 NULL
	// CHARACTER_SET_NAME       string `db:"CHARACTER_SET_NAME"`       // 如果该列是字符数据或 text 数据类型，那么为字符集返回唯一的名称。否则，返回 NULL
	// COLLATION_CATALOG        string `db:"COLLATION_CATALOG"`        // 如果列是字符数据或 text 数据类型，那么返回 master，指明在其中定义排序次序的数据库。否则此列为 NULL
	// COLLATION_SCHEMA         string `db:"COLLATION_SCHEMA"`         // 返回 DBO，为字符数据或 text 数据类型指明排序次序的所有者。否则，返回 NULL
	// COLLATION_NAME           string `db:"COLLATION_NAME"`           // 如果列是字符数据或 text 数据类型，那么为排序次序返回唯一的名称。否则，返回 NULL
	// DOMAIN_CATALOG           string `db:"DOMAIN_CATALOG"`           // 如果列是一种用户定义数据类型，那么该列是某个数据库名称，在该数据库名中创建了这种用户定义数据类型。否则，返回 NULL
	// DOMAIN_SCHEMA            string `db:"DOMAIN_SCHEMA"`            // 如果列是一种用户定义数据类型，那么该列是这种用户定义数据类型的创建者。否则，返回 NULL
	// DOMAIN_NAME              string `db:"DOMAIN_NAME"`              //

	TABLE_CATALOG            string  `db:"TABLE_CATALOG"`            // 表限定符
	TABLE_SCHEMA             string  `db:"TABLE_SCHEMA"`             // 表所有者
	TABLE_NAME               string  `db:"TABLE_NAME"`               // 表名
	COLUMN_NAME              string  `db:"COLUMN_NAME"`              // 列名
	ORDINAL_POSITION         string  `db:"ORDINAL_POSITION"`         // 列标识号
	COLUMN_DEFAULT           *string `db:"COLUMN_DEFAULT"`           // 列的默认值
	IS_NULLABLE              string  `db:"IS_NULLABLE"`              // 列的为空性。如果列允许 NULL，那么该列返回 YES。否则，返回 NO
	DATA_TYPE                string  `db:"DATA_TYPE"`                // 系统提供的数据类型  如 varchar
	CHARACTER_MAXIMUM_LENGTH *string `db:"CHARACTER_MAXIMUM_LENGTH"` // 以字符为单位的最大长度，适于二进制数据、字符数据，或者文本和图像数据。否则，返回 NULL。有关更多信息，请参见数据类型
	CHARACTER_OCTET_LENGTH   *string `db:"CHARACTER_OCTET_LENGTH"`   // 以字节为单位的最大长度，适于二进制数据、字符数据，或者文本和图像数据。否则，返回 NULL
	NUMERIC_PRECISION        *string `db:"NUMERIC_PRECISION"`        // 近似数字数据、精确数字数据、整型数据或货币数据的精度。否则，返回 NULL
	NUMERIC_SCALE            *string `db:"NUMERIC_SCALE"`            // 近似数字数据、精确数字数据、整型数据或货币数据的精度基数。否则，返回 NULL
	DATETIME_PRECISION       *string `db:"DATETIME_PRECISION"`       // 近似数字数据、精确数字数据、整数数据或货币数据的小数位数。否则，返回 NULL
	CHARACTER_SET_NAME       *string `db:"CHARACTER_SET_NAME"`       // datetime 及 SQL-92 interval 数据类型的子类型代码。对于其它数据类型，返回 NULL
	COLLATION_NAME           *string `db:"COLLATION_NAME"`           // 如果列是字符数据或 text 数据类型，那么返回 master，指明字符集所在的数据库。否则，返回 NULL
	COLUMN_TYPE              string  `db:"COLUMN_TYPE"`              // 类型 如 varchar(25)
	COLUMN_KEY               string  `db:"COLUMN_KEY"`               // 如果等于pri，表示是主键
	EXTRA                    string  `db:"EXTRA"`                    // 定义列的时候的其他信息，例如自增，主键
	PRIVILEGES               string  `db:"PRIVILEGES"`               // 返回 DBO，为字符数据或 text 数据类型指明排序次序的所有者。否则，返回 NULL
	COLUMN_COMMENT           string  `db:"COLUMN_COMMENT"`           // 如果列是字符数据或 text 数据类型，那么为排序次序返回唯一的名称。否则，返回 NULL
	GENERATION_EXPRESSION    string  `db:"GENERATION_EXPRESSION"`    // 如果列是一种用户定义数据类型，那么该列是某个数据库名称，在该数据库名中创建了这种用户定义数据类型。否则，返回 NULL
	SRS_ID                   *string `db:"SRS_ID"`                   // 如果列是一种用户定义数据类型，那么该列是这种用户定义数据类型的创建者。否则，返回 NULL
}

var (
	columns     []informationSchemaColumns // 列  information_schema.columns
	columnsOnce sync.Once
)

func InformationSchemaColumns() []informationSchemaColumns {
	columnsOnce.Do(func() {
		strSql := fmt.Sprintf("select * from information_schema.columns where table_schema not in (%s)", ignoreTableSchema)
		db.Select(&columns, strSql)
	})
	return columns
}
