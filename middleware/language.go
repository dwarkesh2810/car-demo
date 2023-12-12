package middleware

import (
	"github.com/beego/beego/v2/server/web/context"
)

func LanguageMiddleware(ctx *context.Context) {
	var lang string
	lang = ctx.Input.Query("lang")
	if len(lang) == 0 {
		lang = ctx.GetCookie("lang")
		if len(lang) != 0 {
			ctx.Input.SetData("lang", lang)
		} else {
			lang = ctx.Input.Header("Accept-Language")
			if len(lang) < 4 {
				ctx.Input.SetData("lang", "en-US")
			} else {
				ctx.Input.SetData("lang", lang[:5])
			}
		}
	} else {
		ctx.Input.SetData("lang", lang)
	}
	ctx.SetCookie("lang", lang)
}
