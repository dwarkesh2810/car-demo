package healthcheck

import (
	_ "car_demo/routers"
	"errors"

	"github.com/beego/beego/v2/client/orm"

	// // _ "github.com/beego/beego/v2/server/web/swagger"

	_ "github.com/lib/pq"
)

type DatabaseCheck struct {
}

func (dc *DatabaseCheck) Check() error {
	if dc.isConnected() {
		return nil
	} else {
		return errors.New("can't connect database")
	}
}

func (dc *DatabaseCheck) isConnected() bool {
	err := orm.RegisterDriver("postgres", orm.DRPostgres)
	if err != nil {
		return false
	}
	err = orm.RegisterDataBase("default", "postgres", "user=root password=1234 dbname=postgres sslmode=disable")
	return err != nil
}
