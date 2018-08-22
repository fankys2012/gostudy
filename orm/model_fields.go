package orm

import (
	"reflect"
	"strings"
)

type fields struct {
	pk        *fieldInfo
	columns   map[string]*fieldInfo
	fields    map[string]*fieldInfo
	fieldsLow map[string]*fieldInfo
	rels      []*fieldInfo
	orders    []string
	dbcols    []string
}

func newFields() *fields {
	f := new(fields)
	f.fields = make(map[string]*fieldInfo)
	f.columns = make(map[string]*fieldInfo)
	f.fieldsLow = make(map[string]*fieldInfo)
	return f
}
func (f *fields) Add(fi *fieldInfo) bool {
	if f.fields[fi.name] == nil && f.columns[fi.name] == nil {
		f.columns[fi.column] = fi
		f.fields[fi.name] = fi
		f.fieldsLow[strings.ToLower(fi.name)] = fi
	} else {
		return false
	}

	f.orders = append(f.orders, fi.column)
	if fi.dbcol {
		f.dbcols = append(f.dbcols, fi.column)
	}
	return true
}

type fieldInfo struct {
	mi         *modelInfo
	fieldIndex []int
	fieldType  int
	dbcol      bool //table column fk and onetoone
	name       string
	column     string
	auto       bool
	pk         bool
	null       bool
	sf         reflect.StructField
	addrValue  reflect.Value
	inModel    bool
	isFielder  bool // implement Fielder interface
}

//get field info by string, name is prior
func (f *fields) GetByAny(name string) (*fieldInfo, bool) {
	if fi, ok := f.fields[name]; ok {
		return fi, ok
	}
	if fi, ok := f.fieldsLow[strings.ToLower(name)]; ok {
		return fi, ok
	}
	if fi, ok := f.columns[name]; ok {
		return fi, ok
	}
	return nil, false
}

func newFieldInfo(mi *modelInfo, field reflect.Value, sf reflect.StructField, mName string) (fi *fieldInfo, err error) {
	var (
		// tag       string
		tagValue  string
		fieldType int
		tags      map[string]string
		attrs     map[string]bool
		addrField reflect.Value
	)
	fi = new(fieldInfo)

	addrField = field
	// CanAddr() 判断能否取址
	if field.CanAddr() && field.Kind() != reflect.Ptr {
		addrField = field.Addr() //获取地址
		//暂时不实现，没搞清楚是个什么结构
		if _, ok := addrField.Interface().(Fielder); ok {
			//...
		}
	}
	//获取表结构的tag属性 sf.Tag.Get() 获取tag
	tags, attrs = parseStructTag(sf.Tag.Get(defaultStructTagName))
	// size := tags["size"]
	switch f := addrField.Interface().(type) {
	case Fielder:
		//...未深入研究
		fieldType = f.FieldType()
	default:
		//rel 类型处理
		tagValue = tags["rel"]
		if tagValue != "" {

		}
		//...
		fieldType, err = getFieldType(addrField)
		if err != nil {
			goto end
		}
	}

	fi.fieldType = fieldType
	fi.name = sf.Name
	fi.column = getColumnName(sf, "")
	fi.sf = sf
	fi.pk = false
	fi.addrValue = addrField
	fi.auto = attrs["auto"]

	//数据库真实字段
	fi.dbcol = true
end:
	if err != nil {
		return nil, err
	}
	return
}

func getColumnName(sf reflect.StructField, col string) string {
	column := col
	if col == "" {
		column = snakeString(sf.Name)
	}
	return column
}
