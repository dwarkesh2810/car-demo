package test

import (
	"car_demo/conf"
	"car_demo/helper"
	"path/filepath"
	"runtime"
	"testing"

	// "github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	// . "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
	conf.LoadEnv("..")
}

// /home/silicon/Dwarkesh/golang/beego/car_demo/app.env
// /home/silicon/Dwarkesh/golang/beego/car_demo/helper/helper_test.go

func TestHashData(t *testing.T) {
	t.Run("Hashdata And VerifyHashedData", func(t *testing.T) {

		data := "bvcjxvbcxkjvbncx"

		result, err := helper.HashData(data)

		if err != nil {
			t.Fatalf("failed to hashed data err := %s", err.Error())
		}

		verify, _ := helper.VerifyHashedData(data, result)

		if !verify {
			t.Fatalf("failed to verify hashed data err := %s", err.Error())
		}
	})
}

func TestSendMail(t *testing.T) {

	t.Run("Send mail", func(t *testing.T) {

		to := "dwarkesh0007@gmail.com"
		subject := "test mail"
		body := "hello, hw r u?"

		sent := helper.SendMail(to, subject, body)

		if !sent {
			t.Fatalf("failed to send mail")
		}
	})
}
