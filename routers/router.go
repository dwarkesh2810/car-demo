// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact dwarkeshpatel.siliconithub@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"car_demo/controllers"
	"car_demo/middleware"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	uc := &controllers.UsersController{}
	mc := &controllers.MsgController{}
	cc := &controllers.Car_masterController{}

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(uc),
			beego.NSRouter("/register", uc, "post:Post"),
			beego.NSRouter("/getall", uc, "get:GetAll"),
			beego.NSRouter("/:id", uc, "get:GetOne;put:Put;delete:Delete"),
			beego.NSRouter("/login", uc, "post:Login"),
			beego.NSRouter("/sendotp", uc, "post:SendOTP"),
			beego.NSRouter("/verifyotp", uc, "post:VerifyOTP"),
			beego.NSRouter("/forgot_password", uc, "post:ForgetPassword"),
			beego.NSRouter("/update", uc, "put:Put"),
		),
		beego.NSNamespace("/sms",
			beego.NSRouter("/otp", mc, "post:SendSMS"),
			beego.NSRouter("/verify", mc, "post:VerifySMS"),
		),
		beego.NSNamespace("/car",
			beego.NSBefore(middleware.Auth),
			beego.NSRouter("/create", cc, "post:Post"),
		),

		beego.NSNamespace("/cars",
			beego.NSRouter("/getall", cc, "get:GetAll"),
			beego.NSRouter("/:id", cc, "get:GetOne;put:Put;delete:Delete"),
		),
	)

	beego.AddNamespace(ns)
}
