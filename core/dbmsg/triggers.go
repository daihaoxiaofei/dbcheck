package dbmsg

import (
	"dbcheck/pkg/db"
	"fmt"
	"sync"
)

type informationSchemaTriggers struct {
	TRIGGER_CATALOG            string `db:"TRIGGER_CATALOG"`
	TRIGGER_SCHEMA             string `db:"TRIGGER_SCHEMA"`
	TRIGGER_NAME               string `db:"TRIGGER_NAME"`
	EVENT_MANIPULATION         string `db:"EVENT_MANIPULATION"`
	EVENT_OBJECT_CATALOG       string `db:"EVENT_OBJECT_CATALOG"`
	EVENT_OBJECT_SCHEMA        string `db:"EVENT_OBJECT_SCHEMA"`
	EVENT_OBJECT_TABLE         string `db:"EVENT_OBJECT_TABLE"`
	ACTION_ORDER               string `db:"ACTION_ORDER"`
	ACTION_CONDITION           string `db:"ACTION_CONDITION"`
	ACTION_STATEMENT           string `db:"ACTION_STATEMENT"`
	ACTION_ORIENTATION         string `db:"ACTION_ORIENTATION"`
	ACTION_TIMING              string `db:"ACTION_TIMING"`
	ACTION_REFERENCE_OLD_TABLE string `db:"ACTION_REFERENCE_OLD_TABLE"`
	ACTION_REFERENCE_NEW_TABLE string `db:"ACTION_REFERENCE_NEW_TABLE"`
	ACTION_REFERENCE_OLD_ROW   string `db:"ACTION_REFERENCE_OLD_ROW"`
	ACTION_REFERENCE_NEW_ROW   string `db:"ACTION_REFERENCE_NEW_ROW"`
	CREATED                    string `db:"CREATED"`
	SQL_MODE                   string `db:"SQL_MODE"`
	DEFINER                    string `db:"DEFINER"`
	CHARACTER_SET_CLIENT       string `db:"CHARACTER_SET_CLIENT"`
	COLLATION_CONNECTION       string `db:"COLLATION_CONNECTION"`
	DATABASE_COLLATION         string `db:"DATABASE_COLLATION"`
}

var (
	triggers     []informationSchemaTriggers // 触发器  INFORMATION_SCHEMA.TRIGGERS
	triggersOnce sync.Once
)

func InformationSchemaTriggers() []informationSchemaTriggers {
	triggersOnce.Do(func() {
		strSql := fmt.Sprintf("select TRIGGER_SCHEMA,TRIGGER_NAME,DEFINER,CREATED "+
			"from information_schema.triggers where TRIGGER_SCHEMA not in (%s)", ignoreTableSchema)
		db.Select(&triggers, strSql)
	})
	return triggers
}
