package controllers_test

import (
	"bytes"
	"car_demo/conf"
	"car_demo/controllers"
	"log"
	"path/filepath"
	"runtime"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/client/orm"
	// "github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
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

func TestCreateUser(t *testing.T) {
	ClearData("users")

	t.Run("CreateUser", func(t *testing.T) {
		Ctrl := &controllers.UsersController{}
		endPoint := "/v1/user/create"

		var jsonStr = []byte(`{"first_name":"Dwarkesh", "last_name":"Patel", "email":"dwarkesh0007@gmail.com", "mobile":"1234324343543", "password":"1234567", "role":"user"}`)

		req, _ := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))

		w := httptest.NewRecorder()

		router := beego.NewControllerRegister()
		router.Add(endPoint, Ctrl)
		router.ServeHTTP(w, req)

		log.Print(w.Body)
		Convey("Subject:Create User Endpoint\n", t, func() {
			Convey("Status Code Should Be 201", func() {
				So(w.Code, ShouldEqual, 201)
			})
			Convey("The Result Should Not Be Empty", func() {
				So(w.Body.Len(), ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestGetAll(t *testing.T) {

	Ctrl := &controllers.UsersController{}
	endPoint := "/v1/user/users"

	req, _ := http.NewRequest("GET", endPoint, nil)

	w := httptest.NewRecorder()
	router := beego.NewControllerRegister()
	router.Add(endPoint, Ctrl)

	router.ServeHTTP(w, req)

	log.Print("1111111111111111111111111111111111111 := Error Code  ", w.Code)

	// Convey("Subject: Get All Endpoint\n", t, func() {
	// 	Convey("Status Code Should Be 200", func() {
	// 		So(w.Code, ShouldEqual, 200)
	// 	})
	// })
	// 	// 	// Convey("The Result Should Not Be Empty", func() {
	// 	// 	// 	So(w.Body.Len(), ShouldBeGreaterThan, 0)
	// 	// 	// })
	// })

}

func TestSendOTP(t *testing.T) {
	TestCreateUser(t)

	t.Run("SendOTP", func(t *testing.T) {
		Ctrl := &controllers.UsersController{}
		endPoints := "/v1/user/sendotp"

		var jsonStrs = []byte(`{"email":"dwarkesh0007@gmail.com"}`)

		req, _ := http.NewRequest("POST", endPoints, bytes.NewBuffer(jsonStrs))

		w := httptest.NewRecorder()

		router := beego.NewControllerRegister()
		router.Add(endPoints, Ctrl)
		router.ServeHTTP(w, req)

		log.Print(w.Code)
		Convey("Subject:Create User Endpoint\n", t, func() {
			Convey("Status Code Should Be 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
			Convey("The Result Should Not Be Empty", func() {
				So(w.Body.Len(), ShouldBeGreaterThan, 0)
			})
		})
	})
}
