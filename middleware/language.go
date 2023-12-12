package middleware

import (
	"github.com/beego/beego/v2/server/web/context"
)

func LanguageMiddleware(ctx *context.Context) {
	lang := ctx.Input.Header("Accept-Language")

	if lang == "" {
		ctx.Input.SetData("lang", "en-US")
	} else {
		ctx.Input.SetData("lang", lang)
	}
}
