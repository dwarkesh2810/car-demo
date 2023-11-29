package main

import (
	"car_demo/conf"
	_ "car_demo/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func init() {
	conf.LoadEnv(".")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=postgres password=1234 dbname=postgres sslmode=disable")
	orm.RunSyncdb("default", false, true)
}

func main() {
	beego.Run()
	// var Data struct {
	// 	Status  string      `json:"status"`
	// 	Data    interface{} `json:"data"`
	// 	Message string      `json:"message"`
	// }

	// // var PostData struct {
	// // 	Name   string `json:"name"`
	// // 	Salary string `json:"salary"`
	// // 	Age    string `json:"age"`
	// // }

	// // JsonData, _ := json.Marshal(&PostData)

	// req, _ := http.NewRequest("Get", "https://dummy.restapiexample.com/api/v1/employees", nil)

	// // req1, _ := http.NewRequest("Post", "https://dummy.restapiexample.com/api/v1/create", bytes.NewBuffer(JsonData))

	// client := &http.Client{}
	// resp, _ := client.Do(req)

	// body, _ := io.ReadAll(resp.Body)

	// json.Unmarshal(body, &Data)

	// log.Print(string(body))
}
