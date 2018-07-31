package orm

import (
	"database/sql"
	"fmt"
	"time"
)

type DrvierType int

const (
	_DriverType = iota
	DRMySQL
)

type alias struct {
	Name string
	// Driver       DriverType
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

	//db ping
	return al, nil
}
