package orm

import (
	"database/sql"
	"reflect"
	"time"
)

type IDbQuerier interface {
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type IOrmer interface {
	Insert(interface{}) (int64, error)
}

type IDbBaser interface {
	Insert(IDbQuerier, *modelInfo, reflect.Value, *time.Location) (int64, error)
}
