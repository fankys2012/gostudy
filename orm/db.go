package orm

import (
	"fmt"
	"reflect"
	"time"
)

// an instance of dbBaser interface/
type dbBase struct {
	ins IDbBaser
}

// execute insert sql dbQuerier with given struct reflect.Value.
func (d *dbBase) Insert(q IDbQuerier, mi *modelInfo, ind reflect.Value, tz *time.Location) (int64, error) {
	// names := make([]string, 0, len(mi.fields.dbcols))
	names := make([]string, 0, 5)
	fmt.Println(names)
	// values, autoFields, err := d.collectValues(mi, ind, mi.fields.dbcols, false, true, &names, tz)
	// if err != nil {
	// 	return 0, err
	// }

	// id, err := d.InsertValue(q, mi, false, names, values)
	// if err != nil {
	// 	return 0, err
	// }

	// if len(autoFields) > 0 {
	// 	err = d.ins.setval(q, mi, autoFields)
	// }
	// return id, err
	return 0, nil
}
