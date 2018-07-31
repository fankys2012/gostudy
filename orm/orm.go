package orm

import (
	"fmt"
	"reflect"
)

type orm struct {
	db IDbQuerier
}

func NewOrm() IOrmer {
	// BootStrap() // execute only once

	o := new(orm)
	// err := o.Using("default")
	// if err != nil {
	// 	panic(err)
	// }
	return o
}

// get model info and model reflect value
func (o *orm) getMiInd(md interface{}, needPtr bool) (mi *modelInfo, ind reflect.Value) {
	val := reflect.ValueOf(md)
	ind = reflect.Indirect(val)
	typ := ind.Type()
	if needPtr && val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<Ormer> cannot use non-ptr model struct `%s`", getFullName(typ)))
	}
	// name := getFullName(typ)
	name := getTableName(val)
	if mi, ok := modelCache.get(name); ok {
		return mi, ind
	}
	panic(fmt.Errorf("<Ormer> table: `%s` not found, make sure it was registered with `RegisterModel()`", name))
}

func (o *orm) Insert(md interface{}) (int64, error) {
	mi, ind := o.getMiInd(md, true)
	// mysql := orm.NewBaseMysql()
	fmt.Println(md)
	return 0, nil
}
