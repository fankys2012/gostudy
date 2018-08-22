package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	formatTime     = "15:04:05"
	formatDate     = "2006-01-02"
	formatDateTime = "2006-01-02 15:04:05"
)

// an instance of dbBaser interface/
type dbBase struct {
	ins IDbBaser
}

// check dbBase implements dbBaser interface.
var _ IDbBaser = new(dbBase)

func (d *dbBase) collectValues(mi *modelInfo, ind reflect.Value, cols []string, names *[]string, insert bool) (values []interface{}, autoFields []string, err error) {
	if names == nil {
		ns := make([]string, 0, len(cols))
		names = &ns
	}
	values = make([]interface{}, 0, len(cols))
	for _, column := range cols {
		var fi *fieldInfo
		if fi, _ = mi.fields.GetByAny(column); fi != nil {
			column = fi.column
		} else {
			panic(fmt.Errorf("wrong db field/column name `%s` for model `%s`", column, mi.fullName))
		}
		//自增
		if fi.auto && insert {
			continue
		}
		value, err := d.collectFieldValues(mi, fi, ind, insert)
		if err != nil {
			return nil, nil, err
		}
		*names = append(*names, column)
		values = append(values, value)
	}
	return
}

func (d *dbBase) collectFieldValues(mi *modelInfo, fi *fieldInfo, ind reflect.Value, insert bool) (interface{}, error) {
	var value interface{}
	if fi.pk {

	} else {
		field := ind.FieldByIndex(fi.fieldIndex)
		if fi.isFielder {

		} else {
			switch fi.fieldType {
			case TypeBooleanField:
				if nb, ok := field.Interface().(sql.NullBool); ok {
					value = nil
					if nb.Valid {
						value = nb.Bool
					}
				} else if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						value = nil
					} else {
						value = field.Elem().Bool()
					}
				} else {
					value = field.Bool()
				}
			case TypeVarCharField, TypeCharField, TypeTextField, TypeJSONField, TypeJsonbField:
				if ns, ok := field.Interface().(sql.NullString); ok {
					value = nil
					if ns.Valid {
						value = ns.String
					}
				} else if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						value = nil
					} else {
						value = field.Elem().String()
					}
				} else {
					value = field.String()
				}
			case TypeSmallIntegerField, TypeIntegerField, TypeBigIntegerField:
				if ni, ok := field.Interface().(sql.NullInt64); ok {
					value = nil
					if ni.Valid {
						value = ni.Int64
					}
				} else if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						value = nil
					} else {
						value = field.Elem().Int()
					}
				} else {
					value = field.Int()
				}
			case TypeTimeField, TypeDateField, TypeDateTimeField:
				value = field.Interface()
				if t, ok := value.(time.Time); ok {
					value = t
				}
			default:

			}
		}
	}
	return value, nil
}

func (d *dbBase) InsertValue(q IDbQuerier, mi *modelInfo, names []string, values []interface{}, isMulti bool) (int64, error) {
	// Q := d.ins.TableQuote()
	Q := "`"
	marks := make([]string, len(names))
	for i := range marks {
		marks[i] = "?"
	}

	sep := fmt.Sprintf("%s,%s", Q, Q)
	qmarks := strings.Join(marks, ", ") // ?,?,? ....
	columns := strings.Join(names, sep) // name`,`age ??

	if isMulti {

	}

	query := fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s) ", Q, mi.table, Q, Q, columns, Q, qmarks)

	fmt.Println(query)
	res, err := q.Exec(query, values...)
	debugLogQueies(query, err, values)
	if err == nil {
		if isMulti {
			return res.RowsAffected()
		}
		return res.LastInsertId()
	}
	return 0, err

}

// execute insert sql dbQuerier with given struct reflect.Value.
func (d *dbBase) Insert(q IDbQuerier, mi *modelInfo, ind reflect.Value, tz *time.Location) (int64, error) {
	names := make([]string, 0, len(mi.fields.dbcols))
	values, autoFields, err := d.collectValues(mi, ind, mi.fields.dbcols, &names, true)
	if err != nil {
		return 0, err
	}

	id, err := d.InsertValue(q, mi, names, values, false)
	if err != nil {
		return 0, err
	}

	if len(autoFields) > 0 {
		// err = d.ins.setval(q, mi, autoFields)
	}
	return id, err
}

func (d *dbBase) Read(q IDbQuerier, mi *modelInfo, ind reflect.Value, cols []string) error {
	var whereCols []string
	var args []interface{}
	if len(cols) > 0 {
		var err error
		whereCols = make([]string, 0, len(cols))
		args, _, err = d.collectValues(mi, ind, cols, &whereCols, false)
		if err != nil {
			return err
		}
	}
	fields := strings.Join(mi.fields.dbcols, "`,`")  // name`,`sex`,`age
	wheres := strings.Join(whereCols, "` = ? AND `") // name`=? and `

	query := fmt.Sprintf("SELECT `%s` FROM `%s` WHERE `%s` = ?", fields, mi.table, wheres)
	//接收数据容器
	refs := make([]interface{}, len(mi.fields.dbcols))
	for i := range refs {
		var ref interface{}
		refs[i] = &ref
	}

	//QueryRow() -查询单条记录
	rows := q.QueryRow(query, args...)

	if err := rows.Scan(refs...); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRows
		}
		return err
	}
	fmt.Println(refs)
	//将获取到的数据塞回model对象内
	// elm := reflect.New(mi.addrField.Elem().Type())
	// mind := reflect.Indirect(elm)
	d.setColsValues(mi, ind, mi.fields.dbcols, refs)
	// d.setColsValues(mi, mind, mi.fields.dbcols, refs)
	// ind.Set(mind)
	return nil
}

func (d *dbBase) setColsValues(mi *modelInfo, ind reflect.Value, cols []string, values []interface{}) {
	for i, column := range cols {
		//获取值
		val := reflect.Indirect(reflect.ValueOf(values[i])).Interface()

		//获取field
		fi := mi.fields.columns[column]
		field := ind.FieldByIndex(fi.fieldIndex)
		value, err := d.convertValueFromDB(fi, val)
		if err != nil {
			panic(fmt.Errorf("Raw value: `%v` %s", val, err.Error()))
		}
		_, err = d.setFieldValue(fi, value, field)

		if err != nil {
			panic(fmt.Errorf("Raw value: `%v` %s", val, err.Error()))
		}
	}
}

// set one value to struct column field.
func (d *dbBase) setFieldValue(fi *fieldInfo, value interface{}, field reflect.Value) (interface{}, error) {

	fieldType := fi.fieldType
	isNative := !fi.isFielder

	switch {
	case fieldType == TypeBooleanField:
		if isNative {
			if nb, ok := field.Interface().(sql.NullBool); ok {
				if value == nil {
					nb.Valid = false
				} else {
					nb.Bool = value.(bool)
					nb.Valid = true
				}
				field.Set(reflect.ValueOf(nb))
			} else if field.Kind() == reflect.Ptr {
				if value != nil {
					v := value.(bool)
					field.Set(reflect.ValueOf(&v))
				}
			} else {
				if value == nil {
					value = false
				}
				field.SetBool(value.(bool))
			}
		}
	case fieldType == TypeVarCharField || fieldType == TypeCharField || fieldType == TypeTextField || fieldType == TypeJSONField || fieldType == TypeJsonbField:
		if isNative {
			if ns, ok := field.Interface().(sql.NullString); ok {
				if value == nil {
					ns.Valid = false
				} else {
					ns.String = value.(string)
					ns.Valid = true
				}
				field.Set(reflect.ValueOf(ns))
			} else if field.Kind() == reflect.Ptr {
				if value != nil {
					v := value.(string)
					field.Set(reflect.ValueOf(&v))
				}
			} else {
				if value == nil {
					value = ""
				}
				field.SetString(value.(string))
			}
		}
	case fieldType == TypeTimeField || fieldType == TypeDateField || fieldType == TypeDateTimeField:
		if isNative {
			if value == nil {
				value = time.Time{}
			} else if field.Kind() == reflect.Ptr {
				if value != nil {
					v := value.(time.Time)
					field.Set(reflect.ValueOf(&v))
				}
			} else {
				field.Set(reflect.ValueOf(value))
			}
		}
	case fieldType == TypePositiveBitField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := uint8(value.(uint64))
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType == TypePositiveSmallIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := uint16(value.(uint64))
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType == TypePositiveIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			if field.Type() == reflect.TypeOf(new(uint)) {
				v := uint(value.(uint64))
				field.Set(reflect.ValueOf(&v))
			} else {
				v := uint32(value.(uint64))
				field.Set(reflect.ValueOf(&v))
			}
		}
	case fieldType == TypePositiveBigIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := value.(uint64)
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType == TypeBitField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := int8(value.(int64))
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType == TypeSmallIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := int16(value.(int64))
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType == TypeIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			if field.Type() == reflect.TypeOf(new(int)) {
				v := int(value.(int64))
				field.Set(reflect.ValueOf(&v))
			} else {
				v := int32(value.(int64))
				field.Set(reflect.ValueOf(&v))
			}
		}
	case fieldType == TypeBigIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := value.(int64)
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType&IsIntegerField > 0:
		if fieldType&IsPositiveIntegerField > 0 {
			if isNative {
				if value == nil {
					value = uint64(0)
				}
				field.SetUint(value.(uint64))
			}
		} else {
			if isNative {
				if ni, ok := field.Interface().(sql.NullInt64); ok {
					if value == nil {
						ni.Valid = false
					} else {
						ni.Int64 = value.(int64)
						ni.Valid = true
					}
					field.Set(reflect.ValueOf(ni))
				} else {
					if value == nil {
						value = int64(0)
					}
					field.SetInt(value.(int64))
				}
			}
		}
	case fieldType == TypeFloatField || fieldType == TypeDecimalField:
		if isNative {
			if nf, ok := field.Interface().(sql.NullFloat64); ok {
				if value == nil {
					nf.Valid = false
				} else {
					nf.Float64 = value.(float64)
					nf.Valid = true
				}
				field.Set(reflect.ValueOf(nf))
			} else if field.Kind() == reflect.Ptr {
				if value != nil {
					if field.Type() == reflect.TypeOf(new(float32)) {
						v := float32(value.(float64))
						field.Set(reflect.ValueOf(&v))
					} else {
						v := value.(float64)
						field.Set(reflect.ValueOf(&v))
					}
				}
			} else {

				if value == nil {
					value = float64(0)
				}
				field.SetFloat(value.(float64))
			}
		}
	}

	if !isNative {
		fd := field.Addr().Interface().(Fielder)
		err := fd.SetRaw(value)
		if err != nil {
			err = fmt.Errorf("converted value `%v` set to Fielder `%s` failed, err: %s", value, fi.name, err)
			return nil, err
		}
	}

	return value, nil
}

// convert value from database result to value following in field type.
func (d *dbBase) convertValueFromDB(fi *fieldInfo, val interface{}) (interface{}, error) {
	if val == nil {
		return nil, nil
	}

	var value interface{}
	var tErr error

	var str *StrTo
	switch v := val.(type) {
	case []byte:
		s := StrTo(string(v))
		str = &s
	case string:
		s := StrTo(v)
		str = &s
	}

	fieldType := fi.fieldType

setValue:
	switch {
	case fieldType == TypeBooleanField:
		if str == nil {
			switch v := val.(type) {
			case int64:
				b := v == 1
				value = b
			default:
				s := StrTo(ToStr(v))
				str = &s
			}
		}
		if str != nil {
			b, err := str.Bool()
			if err != nil {
				tErr = err
				goto end
			}
			value = b
		}
	case fieldType == TypeVarCharField || fieldType == TypeCharField || fieldType == TypeTextField || fieldType == TypeJSONField || fieldType == TypeJsonbField:
		if str == nil {
			value = ToStr(val)
		} else {
			value = str.String()
		}
	case fieldType == TypeTimeField || fieldType == TypeDateField || fieldType == TypeDateTimeField:
		if str == nil {
			switch t := val.(type) {
			case time.Time:
				value = t
			default:
				s := StrTo(ToStr(t))
				str = &s
			}
		}
		if str != nil {
			s := str.String()
			var (
				t   time.Time
				err error
			)
			if len(s) >= 19 {
				s = s[:19]
				// t, err = time.ParseInLocation(formatDateTime, s, tz)
			} else if len(s) >= 10 {
				if len(s) > 10 {
					s = s[:10]
				}
				// t, err = time.ParseInLocation(formatDate, s, tz)
			} else if len(s) >= 8 {
				if len(s) > 8 {
					s = s[:8]
				}
				// t, err = time.ParseInLocation(formatTime, s, tz)
			}
			t = t.In(DefaultTimeLoc)

			if err != nil && s != "00:00:00" && s != "0000-00-00" && s != "0000-00-00 00:00:00" {
				tErr = err
				goto end
			}
			value = t
		}
	case fieldType&IsIntegerField > 0:
		if str == nil {
			s := StrTo(ToStr(val))
			str = &s
		}
		if str != nil {
			var err error
			switch fieldType {
			case TypeBitField:
				_, err = str.Int8()
			case TypeSmallIntegerField:
				_, err = str.Int16()
			case TypeIntegerField:
				_, err = str.Int32()
			case TypeBigIntegerField:
				_, err = str.Int64()
			case TypePositiveBitField:
				_, err = str.Uint8()
			case TypePositiveSmallIntegerField:
				_, err = str.Uint16()
			case TypePositiveIntegerField:
				_, err = str.Uint32()
			case TypePositiveBigIntegerField:
				_, err = str.Uint64()
			}
			if err != nil {
				tErr = err
				goto end
			}
			if fieldType&IsPositiveIntegerField > 0 {
				v, _ := str.Uint64()
				value = v
			} else {
				v, _ := str.Int64()
				value = v
			}
		}
	case fieldType == TypeFloatField || fieldType == TypeDecimalField:
		if str == nil {
			switch v := val.(type) {
			case float64:
				value = v
			default:
				s := StrTo(ToStr(v))
				str = &s
			}
		}
		if str != nil {
			v, err := str.Float64()
			if err != nil {
				tErr = err
				goto end
			}
			value = v
		}
	case fieldType&IsRelField > 0:
		// fi = fi.relModelInfo.fields.pk
		fieldType = fi.fieldType
		goto setValue
	}

end:
	if tErr != nil {
		err := fmt.Errorf("convert to `%s` failed, field: %s err: %s", fi.addrValue.Type(), fi.name, tErr)
		return nil, err
	}

	return value, nil

}

// return quote.
func (d *dbBase) TableQuote() string {
	return "`"
}

// replace value placeholder in parametered sql string.
func (d *dbBase) ReplaceMarks(query *string) {
	// default use `?` as mark, do nothing
}

func debugLogQueies(query string, err error, args ...interface{}) {
	flag := " OK"
	if err != nil {
		flag = "FAIL"
	}
	con := fmt.Sprintf("-[Queries/%s] [%s]", flag, query)
	cons := make([]string, 0, len(args))
	for _, arg := range args {
		cons = append(cons, fmt.Sprintf("%v", arg))
	}
	if len(cons) > 0 {
		con += fmt.Sprintf(" - `%s`", strings.Join(cons, "`, `"))
	}
	if err != nil {
		con += " - " + err.Error()
	}
	fmt.Println(con)
}
