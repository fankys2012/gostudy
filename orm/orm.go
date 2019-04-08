package orm

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// Define common vars
var (
	Debug = false
	// DebugLog         = NewLog(os.Stdout)
	DefaultRowsLimit = 1000
	DefaultRelsDepth = 2
	DefaultTimeLoc   = time.Local
	ErrTxHasBegan    = errors.New("<Ormer.Begin> transaction already begin")
	ErrTxDone        = errors.New("<Ormer.Commit/Rollback> transaction not begin")
	ErrMultiRows     = errors.New("<QuerySeter> return multi rows")
	ErrNoRows        = errors.New("<QuerySeter> no row found")
	ErrStmtClosed    = errors.New("<QuerySeter> stmt already closed")
	ErrArgs          = errors.New("<Ormer> args error may be empty")
	ErrNotImplement  = errors.New("have not implement")
)

type orm struct {
	alias *alias
	db    IDbQuerier
}

func NewOrm() IOrmer {
	// BootStrap() // execute only once

	o := new(orm)
	err := o.Using("default")
	if err != nil {
		panic(err)
	}
	return o
}

//chose db
func (o *orm) Using(name string) error {
	if al, ok := dataBaseCache.get(name); ok {
		o.alias = al
		o.db = al.DB
	} else {
		return fmt.Errorf("orm.Using unknown db name `%s`", name)
	}
	return nil
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
	id, err := o.alias.DbBaser.Insert(o.db, mi, ind, o.alias.TZ)
	fmt.Println(id, err)
	return 0, nil
}

func (o *orm) Read(md interface{}, whereCols []string) error {
	mi, ind := o.getMiInd(md, true)
	err := o.alias.DbBaser.Read(o.db, mi, ind, whereCols)
	return err
}

func (o *orm) Update(md interface{}, whereCols []string) (int64, error) {

	return 0, nil
}

func (o *orm) Delete(md interface{}, whereCols []string) (int64, error) {
	return 0, nil
}
