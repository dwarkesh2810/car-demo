package test

import (
	"bytes"
	"car_demo/conf"
	"car_demo/helper"
	"car_demo/middleware"
	"car_demo/request"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"

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

func TestRouters(Ctrl beego.ControllerInterface, endPoint, token, method, mappedMethod string, jsonString []byte, isAuthRequired bool) *httptest.ResponseRecorder {
	var jwt string
	var cw *httptest.ResponseRecorder
	var router *beego.ControllerRegister
	req, _ := http.NewRequest(method, endPoint, bytes.NewBuffer(jsonString))

	if isAuthRequired {
		jwt = fmt.Sprintf("Bearer %s", token)
		req.Header.Set("Authorization", jwt)
		cw = httptest.NewRecorder()
		router = beego.NewControllerRegister()
		router.InsertFilter(endPoint, beego.BeforeRouter, middleware.Auth, beego.WithCaseSensitive(true))
	} else {
		cw = httptest.NewRecorder()
		router = beego.NewControllerRegister()
	}

	router.Add(endPoint, Ctrl, beego.WithRouterMethods(Ctrl, fmt.Sprintf("%s:%s", strings.ToLower(method), mappedMethod)))

	router.ServeHTTP(cw, req)
	return cw
}

func GetresponseDate(body *bytes.Buffer) (map[string]interface{}, error) {
	var u helper.Response
	b, err := io.ReadAll(body)

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, err
	}
	resp := map[int]interface{}{1: u.Data}
	return resp[1].(map[string]interface{}), nil
}
