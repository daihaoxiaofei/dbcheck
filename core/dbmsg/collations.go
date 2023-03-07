package dbmsg

import (
	"dbcheck/pkg/db"
	"sync"
)

//  information_schema.COLLATIONS 关于各字符集的对照信息
type informationSchemaCollations struct {
	COLLATION_NAME     string `db:"COLLATION_NAME"`     // 排序规则名  utf8mb4_general_ci
	CHARACTER_SET_NAME string `db:"CHARACTER_SET_NAME"` // 这个排序规则所对应的字符集名字 utf8mb4
	ID                 string `db:"ID"`                 // 排序规则的ID
	IS_DEFAULT         string `db:"IS_DEFAULT"`         // 是否是字符集的默认排序规则
	IS_COMPILED        string `db:"IS_COMPILED"`        // 是不已经被编译进MySQL来了
	SORTLEN            string `db:"SORTLEN"`            //
	PAD_ATTRIBUTE      string `db:"PAD_ATTRIBUTE"`      //
}

// armscii8_general_ci	armscii8	32		Yes	Yes	1	PAD SPACE
// armscii8_bin			armscii8	64		Yes	1	PAD SPACE
// ascii_general_ci		ascii		11		Yes	Yes	1	PAD SPACE
// ascii_bin			ascii		65		Yes	1	PAD SPACE

// var (
// 	collations     []informationSchemaCollations // information_schema.COLLATIONS
// 	collationsOnce sync.Once
// )
//
// func InformationSchemaCollations() []informationSchemaCollations {
// 	collationsOnce.Do(func() {
// 		db.Select(&collations, "select * from information_schema.COLLATIONS")
// 	})
// 	return collations
// }

var (
	collations     map[string]string // information_schema.COLLATIONS
	collationsOnce sync.Once         // information_schema.COLLATIONS
)

func CollationsMap() map[string]string {
	collationsOnce.Do(func() {
		collations = db.GetMap(`select COLLATION_NAME, CHARACTER_SET_NAME from information_schema.COLLATIONS`)
	})
	return collations

}
