package main

import (
	"car_demo/conf"
	"car_demo/logger"
	_ "car_demo/routers"
	"car_demo/task"
	"context"
	"log"

	"github.com/beego/beego/v2/client/orm"

	// // _ "github.com/beego/beego/v2/server/web/swagger"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func init() {
	conf.LoadEnv(".")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=root password=1234 dbname=postgres sslmode=disable")
	orm.RunSyncdb("default", false, true)
	logger.Init()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	task.CreateTask("test1", "0 42 15 * * *", Demo)
	beego.Run()
}

func Demo(c context.Context) error {
	log.Print("hello")
	return nil
}
