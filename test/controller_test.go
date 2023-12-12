package test

import (
	"car_demo/controllers"

	"net/http"
	"testing"

	// "github.com/beego/beego/v2/core/logs"

	_ "github.com/beego/beego/v2/server/web/context"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateUser(t *testing.T) {
	ClearData("users")

	t.Run("CreateUser", func(t *testing.T) {
		Ctrl := &controllers.UsersController{}
		endPoint := "/v1/users/register"
		var jsonStr = []byte(`{"first_name":"Dwarkesh", "last_name":"Patel", "email":"dwarkesh0007@gmail.com", "mobile":"12343241343543", "password":"1234567", "role":"user"}`)

		token := ""
		method := "POST"
		mappedMethod := "Register"

		w := TestRouters(Ctrl, endPoint, token, method, mappedMethod, jsonStr, false)

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
		endPoint := "/v1/user/getall"
		mappedMethod := "GetAll"
		token := ""
		method := "GET"

		w := TestRouters(Ctrl, endPoint, token, method, mappedMethod, nil, false)

		Convey("Subject: Get All Users\n", t, func() {
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
	AddUser()
	t.Run("SendOTP", func(t *testing.T) {

		Ctrl := &controllers.UsersController{}
		endPoints := "/v1/users/sendotp"
		var jsonStr = []byte(`{"email":"dwarkesh0007@gmail.com"}`)
		token := ""
		method := "POST"
		mappedMethod := "SendOTP"

		w := TestRouters(Ctrl, endPoints, token, method, mappedMethod, jsonStr, false)

		Convey("Subject:Send OTP\n", t, func() {
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
		var jsonStr = []byte(`{"email":"dwarkesh0007@gmail.com","otp":"46314"}`)
		token := ""
		method := "POST"
		mappedMethod := "VerifyOTP"

		w := TestRouters(Ctrl, endPoints, token, method, mappedMethod, jsonStr, false)

		Convey("Subject:Verify OTP\n", t, func() {
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
		var jsonStr = []byte(`{"email":"dwarkesh00071@gmail.com","password":"1234567"}`)
		token := ""
		method := "POST"
		mappedMethod := "Login"

		w := TestRouters(Ctrl, endPoints, token, method, mappedMethod, jsonStr, false)

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
		jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE0MjgwMDMsInN1YiI6MjA4fQ.0qThD22k8IxZf4Fb65o-FHdgzwK5I55_Z8LaVDK-VCI"
		var jsonStr = []byte(`{"car_name":"BMW X7", "make":"BMW", "model":"X7", "car_type":"SUV"}`)
		method := "POST"
		mappedMethod := "Post"
		result := TestRouters(Ctrl, endPoint, jwt, method, mappedMethod, jsonStr, true)
		Convey("Subject:Create User Endpoint\n", t, func() {
			Convey("Status Code Should Be 200", func() {
				So(result.Code, ShouldEqual, http.StatusCreated)
			})
		})
	})
}
