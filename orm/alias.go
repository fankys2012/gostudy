package orm

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type DriverType int

const (
	_DriverType = iota
	DRMySQL
)

type dbCache struct {
	mux   sync.RWMutex
	cache map[string]*alias
}

func (cache *dbCache) add(name string, al *alias) bool {
	cache.mux.Lock()
	defer cache.mux.Unlock()
	if _, ok := cache.cache[name]; !ok {
		cache.cache[name] = al
		return true
	}
	return false
}

func (cache *dbCache) get(name string) (al *alias, ok bool) {
	cache.mux.Lock()
	defer cache.mux.Unlock()
	al, ok = cache.cache[name]
	return
}

var (
	drivers = map[string]DriverType{
		"mysql": DRMySQL,
	}
	dataBaseCache = &dbCache{cache: make(map[string]*alias)}
)

type alias struct {
	Name         string
	Driver       DriverType
	DriverName   string
	DataSource   string
	MaxIdleConns int
	MaxOpenConns int
	DB           *sql.DB
	DbBaser      IDbBaser
	TZ           *time.Location
	Engine       string
}

func RegisterDataBase(aliasName, driverName, dataSource string, params ...interface{}) error {

	var (
		err error
		db  *sql.DB
		al  *alias
	)

	db, err = sql.Open(driverName, dataSource)
	if err != nil {
		err = fmt.Errorf("Register db %s %s", aliasName, err)
		goto end
	}
	al, err = addAliasWthDB(aliasName, driverName, db)
	al.DataSource = dataSource

end:
	if err != nil {
		if db != nil {
			db.Close()
		}

	}
	return err
}

func addAliasWthDB(aliasName, driverName string, db *sql.DB) (*alias, error) {
	al := new(alias)
	al.Name = aliasName
	al.DriverName = driverName
	al.DB = db
	if dr, ok := drivers[driverName]; ok {
		switch dr {
		case DRMySQL:
			al.DbBaser = NewBaseMysql()
		}
		al.Driver = dr
	}
	//db ping
	if !dataBaseCache.add(aliasName, al) {
		return nil, fmt.Errorf("DataBase alias name `%s` already registered, cannot reuse", aliasName)
	}
	return al, nil
}
