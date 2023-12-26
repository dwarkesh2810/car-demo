package main

import (
	"car_demo/conf"
	"car_demo/healthcheck"
	"car_demo/validations"

	"car_demo/logger"
	_ "car_demo/routers"
	"car_demo/task"

	"context"
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/admin"

	// // _ "github.com/beego/beego/v2/server/web/swagger"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/session/postgres"
	_ "github.com/lib/pq"
)

func init() {
	conf.LoadEnv(".")

	orm.RegisterDriver("postgres", orm.DRPostgres)

	orm.RegisterDataBase("default", "postgres", "user=root password=1234 dbname=postgres host=postgres_db sslmode=disable")
	// orm.RegisterDataBase("default", "postgres", "user=root password=1234 dbname=golang_demo sslmode=disable")
	// beego.BConfig.WebConfig.Session.SessionOn = true
	// beego.BConfig.WebConfig.Session.SessionProvider = "postgresql"
	// beego.BConfig.WebConfig.Session.SessionProviderConfig = "postgres://root:1234@localhost/postgres?sslmode=disable"

	orm.RunSyncdb("default", false, true)
	logger.Init()
	validations.Init()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	// helper.TotalRecords("users")
	admin.AddHealthCheck("database", &healthcheck.DatabaseCheck{})
	task.CreateTask("test1", "0 */1 * * * *", Demo)
	beego.Run()
}

func Demo(c context.Context) error {
	// imports.Seed(25)
	log.Print("hello")
	return nil
}
