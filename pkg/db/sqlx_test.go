package db

import (
	"fmt"
	"testing"
)

type GlobalVariables struct {
	VariableName string `db:"Variable_name"`
	Value        string `db:"Value"`
}

func Test_(t *testing.T) {
	var gv []GlobalVariables
	err := DB.Select(&gv, `show global variables`)
	if err != nil {
		panic(err)
	}

	fmt.Println(gv)
}

func TestGetMap(t *testing.T) {
	m := GetMap(`show global variables`)
	fmt.Println(m)
}
