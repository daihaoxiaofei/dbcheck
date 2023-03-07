package dbmsg

import (
	"dbcheck/pkg/db"
	"sync"
)

// var (
// 	GlobalVariablesMap map[string]string // global variables
// 	GlobalStatusMap                      // global status
// )

var (
	variables     map[string]string // information_schema.COLLATIONS
	variablesOnce sync.Once         // information_schema.COLLATIONS

	status     map[string]string
	statusOnce sync.Once
)

func GlobalVariablesMap() map[string]string {
	variablesOnce.Do(func() {
		variables = db.GetMap(`show global variables`)
	})
	return variables

}

func GlobalStatusMap() map[string]string {
	statusOnce.Do(func() {
		status = db.GetMap(`show global status`)
	})
	return status
}
