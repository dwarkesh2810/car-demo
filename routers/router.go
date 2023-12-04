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
	// mc := &controllers.MsgController{}
	cc := &controllers.Car_masterController{}

	ns := beego.NewNamespace("/v1",
		beego.NSAutoRouter(uc),
		// beego.NSAutoRouter(mc),
		beego.NSInclude(uc),
		// beego.NSNamespace("/sms",
			// beego.NSRouter("/otp", mc, "post:SendSMS"),
		// 	beego.NSRouter("/verify", mc, "post:VerifySMS"),
		// ),
		// beego.NSNamespace("/car",
		// 	beego.NSBefore(middleware.Auth),
		// 	beego.NSRouter("/create", cc, "post:Post"),
		// ),

		// beego.NSNamespace("/cars",
		// 	beego.NSRouter("/getall", cc, "get:GetAll"),
		// 	beego.NSRouter("/:id", cc, "get:GetOne;put:Put;delete:Delete"),
		// ),
	)

	n1 := beego.NewNamespace("/v1",
		beego.NSAutoRouter(cc),
		beego.NSInclude(cc),
		// beego.NSBefore(middleware.Auth),
	)

	beego.InsertFilter("/v1/car_master/create", beego.BeforeRouter, middleware.Auth)
	// beego.Router("v1/users/getone/?:id", uc, "get:GetOne")

	beego.AddNamespace(ns)
	beego.AddNamespace(n1)
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE3NTYzMTUsInN1YiI6MjIwfQ.Ulw6ZLo9mNX6tlrq4yNJ118CsuoUoY1eDZhJPPHb_qo

// {
// 	"car_image": "assets/img/cars/sedan/HuracanModelImage.jpg",
// 	"car_name": "BMW",
// 	"car_type": "sedan",
// 	"make": "BMW",
// 	"model": "X7"
//   }
