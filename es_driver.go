package es_driver

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type ESDriver struct {
}

func (esd *ESDriver) Open(dsn string) (driver.Conn, error) {

	var conf ESConnConfig

	err := json.Unmarshal([]byte(dsn), &conf)
	if err != nil {
		return nil, err
	}

	ec := new(ESConn)
	ec.Init(&conf)

	return ec, nil
}

func init() {
	sql.Register("es", &ESDriver{})
}
