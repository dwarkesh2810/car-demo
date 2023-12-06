package main

import (
	"car_demo/conf"
	"car_demo/export"
	"car_demo/logger"
	"car_demo/models"
	_ "car_demo/routers"
	"car_demo/sessions"
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
	orm.RegisterDataBase("default", "postgres", "user=postgres password=1234 dbname=postgres sslmode=disable")
	orm.RunSyncdb("default", false, true)
	logger.Init()
	sessions.SetSessionName("dexter")
	sessions.SetSessionCookieLifeTime(60)
	sessions.SetSessionGCMaxLifetime(60)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	task.CreateTask("test", "21 * * * * * ", Demo) //sec min hour day month weekday
	var user models.Users
	columns := []string{"Id", "FirstName", "LastName", "Email", "Mobile", "CreatedAt"}

	// export.ExportToCSV(user, columns)
	export.ExportToExcel(user, columns)
	export.DbToPdf(user, columns)
	beego.Run()

	// field := val.FieldByName("ID")
}

func Demo(c context.Context) error {
	log.Print("Hello")
	return nil

}
