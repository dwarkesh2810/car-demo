package test

import (
	"bytes"
	"car_demo/controllers"
	"car_demo/middleware"
	"fmt"
	"log"

	"net/http"
	"net/http/httptest"
	"testing"

	// "github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/context"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateUser(t *testing.T) {
	ClearData("users")

	t.Run("CreateUser", func(t *testing.T) {
		Ctrl := &controllers.UsersController{}
		endPoint := "/v1/user/create"

		var jsonStr = []byte(`{"first_name":"Dwarkesh", "last_name":"Patel", "email":"dwarkesh00071@gmail.com", "mobile":"12343241343543", "password":"1234567", "role":"user"}`)

		req, _ := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))

		w := httptest.NewRecorder()

		router := beego.NewControllerRegister()
		router.Add(endPoint, Ctrl, beego.WithRouterMethods(Ctrl, "post:Post"))
		router.ServeHTTP(w, req)

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

	t.Run("GetAll", func(t *testing.T) {
		Ctrl := &controllers.UsersController{}
		endPoint := "/v1/user/"

		req, _ := http.NewRequest("GET", endPoint, nil)

		w := httptest.NewRecorder()
		router := beego.NewControllerRegister()

		router.Add(endPoint, Ctrl, beego.WithRouterMethods(Ctrl, "get:GetAll"))

		router.ServeHTTP(w, req)

		Convey("Subject: Get All Endpoint\n", t, func() {
			Convey("Status Code Should Be 200", func() {
				So(w.Code, ShouldEqual, 200)
			})

			Convey("The Result Should Not Be Empty", func() {
				So(w.Body.Len(), ShouldBeGreaterThan, 0)
			})
		})
	})

}

func TestSendOTP(t *testing.T) {

	t.Run("SendOTP", func(t *testing.T) {
		Ctrl := &controllers.UsersController{}
		endPoints := "/v1/user/sendotp"

		var jsonStrs = []byte(`{"email":"dwarkesh0007@gmail.com"}`)

		req, _ := http.NewRequest("POST", endPoints, bytes.NewBuffer(jsonStrs))

		w := httptest.NewRecorder()

		router := beego.NewControllerRegister()
		router.Add(endPoints, Ctrl, beego.WithRouterMethods(Ctrl, "post:SendOTP"))
		router.ServeHTTP(w, req)

		// json.Unmarshal()
		fmt.Printf("type of body %T", w.Body)
		Convey("Subject:Create User Endpoint\n", t, func() {
			Convey("Status Code Should Be 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestVerifyOtp(t *testing.T) {
	t.Run("VerifyOtp", func(t *testing.T) {
		Ctrl := &controllers.UsersController{}
		endPoints := "/v1/user/verifyotp"
		var jsonStrs = []byte(`{"email":"dwarkesh0007@gmail.com","otp":"310376"}`)

		req, _ := http.NewRequest("POST", endPoints, bytes.NewBuffer(jsonStrs))

		w := httptest.NewRecorder()

		router := beego.NewControllerRegister()
		router.Add(endPoints, Ctrl, beego.WithRouterMethods(Ctrl, "post:VerifyOTP"))
		router.ServeHTTP(w, req)

		// json.Unmarshal()
		fmt.Printf("type of body %T", w.Body)
		Convey("Subject:Create User Endpoint\n", t, func() {
			Convey("Status Code Should Be 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestLogin(t *testing.T) {
	t.Run("Login", func(t *testing.T) {

		Ctrl := &controllers.UsersController{}
		endPoints := "/v1/user/login"
		var jsonStrs = []byte(`{"email":"dwarkesh0007@gmail.com","password":"1234567"}`)

		req, _ := http.NewRequest("POST", endPoints, bytes.NewBuffer(jsonStrs))

		w := httptest.NewRecorder()

		router := beego.NewControllerRegister()
		router.Add(endPoints, Ctrl, beego.WithRouterMethods(Ctrl, "post:Login"))
		router.ServeHTTP(w, req)

		Convey("Subject:Create User Endpoint\n", t, func() {
			Convey("Status Code Should Be 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestAddCar(t *testing.T) {

	t.Run("AddCar", func(t *testing.T) {

		Ctrl := &controllers.Car_masterController{}

		endPoint := "/v1/cars/create"
		jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE0MjYzOTQsInN1YiI6MjA4fQ.yvv6lZuMwMVCVlLCHCVJaZieeYpQlIYlF90BPp59rzQ"
		token := fmt.Sprintf("Bearer %s", jwt)

		var jsonStr = []byte(`{"car_name":"BMW X7", "make":"BMW", "model":"X7", "car_type":"sedan"}`)

		creq, _ := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		creq.Header.Set("Authorization", token)

		log.Print(creq.Header.Get("Authorization"))
		cw := httptest.NewRecorder()

		crouter := beego.NewControllerRegister()
		crouter.InsertFilter(endPoint, beego.BeforeRouter, middleware.Auth, beego.WithCaseSensitive(true))
		crouter.Add(endPoint, Ctrl, beego.WithRouterMethods(Ctrl, "post:Post"))
		crouter.ServeHTTP(cw, creq)

		Convey("Subject:Create User Endpoint\n", t, func() {
			Convey("Status Code Should Be 200", func() {
				So(cw.Code, ShouldEqual, http.StatusCreated)
			})
		})
	})
}
