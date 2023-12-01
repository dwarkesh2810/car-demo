package main

import (
	"car_demo/conf"
	_ "car_demo/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	// _ "github.com/beego/beego/v2/server/web/swagger"
	_ "github.com/lib/pq"
)

func init() {
	conf.LoadEnv(".")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=postgres password=1234 dbname=postgres sslmode=disable")
	orm.RunSyncdb("default", false, true)
}

func main() {

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Run()

}

// func RegisterRoutes(ns *beego.Namespace, controller *beego.Controller, basePath string) {
// 	ctrlType := reflect.TypeOf(controller)

// 	for i := 0; i < ctrlType.NumMethod(); i++ {
// 		methodName := ctrlType.Method(i).Name

// 		if strings.HasPrefix(methodName, "Mapping") {
// 			handler, _ := ctrlType.MethodByName(methodName)
// 			method := strings.TrimPrefix(methodName, "Mapping")

// 			beego.AutoRouter(controller)
// 		}
// 	}
// }
