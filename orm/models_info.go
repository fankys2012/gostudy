package orm

import (
	"fmt"
	"os"
	"reflect"
)

const (
	defaultStructTagName = "orm" //表结构体中的tag属性名
)

// single model info
type modelInfo struct {
	pkg       string
	name      string
	fullName  string
	table     string
	model     interface{}
	fields    *fields
	manual    bool
	addrField reflect.Value //store the original struct value
	uniques   []string
	isThrough bool
}

// new model info
func newModelInfo(val reflect.Value) (mi *modelInfo) {
	mi = &modelInfo{}
	mi.fields = newFields()
	ind := reflect.Indirect(val)
	mi.addrField = val
	mi.name = ind.Type().Name() // mi.fullName = getFullName(ind.Type())
	addModelFields(mi, ind, "", []int{})
	return
}

func RegisterModel(prefix string, model interface{}) {
	val := reflect.ValueOf(model)
	// typ := reflect.Indirect(val).Type()
	table := getTableName(val)

	if prefix != "" {
		table = prefix + table
	}

	//register model info
	mi := newModelInfo(val)
	mi.table = table
	mi.model = model
	modelCache.set(table, mi)
}

func getTableName(val reflect.Value) string {
	if fun := val.MethodByName("TableName"); fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 && vals[0].Kind() == reflect.String {
			return vals[0].String()
		}
	}
	return snakeString(reflect.Indirect(val).Type().Name())
}

// get reflect.Type name with package path.
func getFullName(typ reflect.Type) string {
	return typ.PkgPath() + "." + typ.Name()
}

func addModelFields(mi *modelInfo, ind reflect.Value, name string, index []int) {
	var (
		err error
		fi  *fieldInfo
		sf  reflect.StructField
	)

	for i := 0; i < ind.NumField(); i++ {
		field := ind.Field(i)
		sf = ind.Type().Field(i)

		fi, err = newFieldInfo(mi, field, sf, name)
		if err != nil {
			break
		}
		fi.fieldIndex = append(index, i)
		fi.mi = mi
		fi.inModel = true
		if !mi.fields.Add(fi) {
			err = fmt.Errorf("duplicate column name: %s", fi.name)
		}
	}

	if err != nil {
		fmt.Println(fmt.Errorf("field :%s.%s,%s", ind.Type(), sf.Name, err))
		os.Exit(2)
	}
}
