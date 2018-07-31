package orm

import "reflect"

// single model info
type modelInfo struct {
	pkg      string
	name     string
	fullName string
	table    string
	model    interface{}
	// fields    *fields
	manual    bool
	addrField reflect.Value //store the original struct value
	uniques   []string
	isThrough bool
}

// new model info
func newModelInfo(val reflect.Value) (mi *modelInfo) {
	mi = &modelInfo{}
	// mi.fields = newFields()
	ind := reflect.Indirect(val)
	mi.addrField = val
	mi.name = ind.Type().Name()
	// mi.fullName = getFullName(ind.Type())
	// addModelFields(mi, ind, "", []int{})
	return
}

func RegisterModel(model interface{}) {
	val := reflect.ValueOf(model)
	// typ := reflect.Indirect(val).Type()
	table := getTableName(val)

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
