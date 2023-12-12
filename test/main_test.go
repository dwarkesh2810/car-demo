package test

import (
	"car_demo/conf"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

func TestMain(t *testing.M) {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
	conf.LoadEnv("..")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=root password=1234 dbname=postgres sslmode=disable")
	orm.RunSyncdb("default", false, true)

	// ClearData("car_master")
	// ClearData("users")
	code := t.Run()
	os.Exit(code)
}
