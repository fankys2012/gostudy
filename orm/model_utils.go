package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var supportTag = map[string]int{
	"size":    1,
	"column":  1,
	"default": 1,
	"rel":     1,
	"reverse": 1,
	"type":    1,
	"auto":    2,
	"pk":      2,
}

func parseStructTag(data string) (tags map[string]string, attrs map[string]bool) {

	attrs = make(map[string]bool)

	tags = make(map[string]string)
	for _, v := range strings.Split(data, ";") {
		if v == "" {
			continue
		}
		v = strings.TrimSpace(v)
		if t := strings.ToLower(v); supportTag[t] == 2 {
			attrs[t] = true
		} else if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
			name := t[:i] //括号之前部分
			if supportTag[name] == 1 {
				v = v[i+1 : len(v)-1] //括号内内容
				tags[name] = v
			}
		}
	}
	return
}

// return field type as type constant from reflect.Value
func getFieldType(val reflect.Value) (ft int, err error) {
	switch val.Type() {
	case reflect.TypeOf(new(int8)):
		ft = TypeBitField
	case reflect.TypeOf(new(int16)):
		ft = TypeSmallIntegerField
	case reflect.TypeOf(new(int32)),
		reflect.TypeOf(new(int)):
		ft = TypeIntegerField
	case reflect.TypeOf(new(int64)):
		ft = TypeBigIntegerField
	case reflect.TypeOf(new(uint8)):
		ft = TypePositiveBitField
	case reflect.TypeOf(new(uint16)):
		ft = TypePositiveSmallIntegerField
	case reflect.TypeOf(new(uint32)),
		reflect.TypeOf(new(uint)):
		ft = TypePositiveIntegerField
	case reflect.TypeOf(new(uint64)):
		ft = TypePositiveBigIntegerField
	case reflect.TypeOf(new(float32)),
		reflect.TypeOf(new(float64)):
		ft = TypeFloatField
	case reflect.TypeOf(new(bool)):
		ft = TypeBooleanField
	case reflect.TypeOf(new(string)):
		ft = TypeVarCharField
	case reflect.TypeOf(new(time.Time)):
		ft = TypeDateTimeField
	default:
		elm := reflect.Indirect(val)
		switch elm.Kind() {
		case reflect.Int8:
			ft = TypeBitField
		case reflect.Int16:
			ft = TypeSmallIntegerField
		case reflect.Int32, reflect.Int:
			ft = TypeIntegerField
		case reflect.Int64:
			ft = TypeBigIntegerField
		case reflect.Uint8:
			ft = TypePositiveBitField
		case reflect.Uint16:
			ft = TypePositiveSmallIntegerField
		case reflect.Uint32, reflect.Uint:
			ft = TypePositiveIntegerField
		case reflect.Uint64:
			ft = TypePositiveBigIntegerField
		case reflect.Float32, reflect.Float64:
			ft = TypeFloatField
		case reflect.Bool:
			ft = TypeBooleanField
		case reflect.String:
			ft = TypeVarCharField
		default:
			if elm.Interface() == nil {
				panic(fmt.Errorf("%s is nil pointer, may be miss setting tag", val))
			}
			switch elm.Interface().(type) {
			case sql.NullInt64:
				ft = TypeBigIntegerField
			case sql.NullFloat64:
				ft = TypeFloatField
			case sql.NullBool:
				ft = TypeBooleanField
			case sql.NullString:
				ft = TypeVarCharField
			case time.Time:
				ft = TypeDateTimeField
			}
		}
	}
	if ft&IsFieldType == 0 {
		err = fmt.Errorf("unsupport field type %s, may be miss setting tag", val)
	}
	return
}

func getPkColumnInfo(mi *modelInfo, ind reflect.Value) (column string, value interface{}, exists bool) {
	fi := mi.fields.pk
	v := ind.FieldByIndex(fi.fieldIndex)

	switch fi.fieldType {
	case TypeVarCharField, TypeCharField, TypeTextField, TypeJSONField, TypeJsonbField:
		value = v.String()
		exists = true
	case TypeSmallIntegerField, TypeIntegerField, TypeBigIntegerField:
		value = v.Int()
		exists = true
	}
	column = fi.column
	return
}

// get pk column info.
func getExistPk(mi *modelInfo, ind reflect.Value) (column string, value interface{}, exist bool) {
	fi := mi.fields.pk

	v := ind.FieldByIndex(fi.fieldIndex)
	if fi.fieldType&IsPositiveIntegerField > 0 {
		vu := v.Uint()
		exist = vu > 0
		value = vu
	} else if fi.fieldType&IsIntegerField > 0 {
		vu := v.Int()
		exist = true
		value = vu
	} else if fi.fieldType&IsRelField > 0 {
		// _, value, exist = getExistPk(fi.relModelInfo, reflect.Indirect(v))
	} else {
		vu := v.String()
		exist = vu != ""
		value = vu
	}

	column = fi.column
	return
}
