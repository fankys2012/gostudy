package orm

import (
	"database/sql"
	"reflect"
	"time"
)

//sql.db 自动实现
type IDbQuerier interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type IOrmer interface {
	Insert(interface{}) (int64, error)
	Read(interface{}, []string) error
	Delete(interface{}, []string) (int64, error)
	Update(interface{}, []string) (int64, error)
}

type IDbBaser interface {
	Insert(IDbQuerier, *modelInfo, reflect.Value, *time.Location) (int64, error)
	Read(IDbQuerier, *modelInfo, reflect.Value, []string) error
	TableQuote() string
	ReplaceMarks(*string)
}

// Fielder define field info
type Fielder interface {
	String() string
	FieldType() int
	SetRaw(interface{}) error
	RawValue() interface{}
}
