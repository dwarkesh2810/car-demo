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

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

func init() {

	langs := []string{"en-US", "zh-CN", "hi-IN"} // List of supported languages

	for _, lang := range langs {
		if err := i18n.SetMessage(lang, "conf/locale_"+lang+".ini"); err != nil {
			// logger.Error("Fail to set message file:", err)
			logs.Error("Fail to set message file:", err)
			return
		}
	}

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

	beego.Router("/demoset", uc, "post:DemoSet")
	beego.Router("/demoget", uc, "get:DemoGet")

	beego.InsertFilter("/v1/car_master/create", beego.BeforeRouter, middleware.Auth)
	beego.InsertFilter("*", beego.BeforeRouter, middleware.LanguageMiddleware)
	// beego.Router("v1/users/getone/?:id", uc, "get:GetOne")

	beego.AddNamespace(ns)
	beego.AddNamespace(n1)
}

// beei18n sync locale_en-US.ini locale_zh-CN.ini
