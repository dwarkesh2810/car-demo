package healthcheck

import (
	_ "car_demo/routers"
	"errors"

	"github.com/beego/beego/v2/client/orm"
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
	o := orm.NewOrm()
	if _, err := o.Raw("SELECT 1").Exec(); err != nil {
		return false
	}
	return true
}
