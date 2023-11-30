package test

import (
	"car_demo/conf"
	"car_demo/request"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/beego/beego/v2/client/orm"

	// "github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func init() {
	conf.LoadEnv("..")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=postgres password=1234 dbname=postgres sslmode=disable")
	orm.RunSyncdb("default", false, true)

	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

}

func ClearData(tableName string) error {
	o := orm.NewOrm()
	// Construct the SQL query to truncate the table
	query := fmt.Sprintf("TRUNCATE TABLE %s", tableName)

	// Execute the raw SQL query to truncate the table
	_, err := o.Raw(query).Exec()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func UserCreateData() *request.CreateUserRequest {
	return &request.CreateUserRequest{
		FirstName: "Dexter",
		LastName:  "Pat",
		Email:     "dwarkesh01007@gmail.com",
		Mobile:    "1235465415656",
		Password:  "123456",
		Role:      "user",
	}

}
