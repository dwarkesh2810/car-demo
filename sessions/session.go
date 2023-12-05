package sessions

import (
	beego "github.com/beego/beego/v2/server/web"
)

func Set(c beego.Controller, key string, value interface{}) error {
	return c.SetSession(key, value)
}

func Get(c beego.Controller, key string) interface{} {
	return c.GetSession(key)
}

func Destroy(c beego.Controller) error {
	return c.DestroySession()
}

func SetKeyToNil(c beego.Controller, key string) error {
	return Set(c, key, nil)
}

func SetSessionName(name string) {
	beego.BConfig.WebConfig.Session.SessionName = name
}

func SetSessionGCMaxLifetime(duration int64) {
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = duration
}

func SetSessionCookieLifeTime(duration int) {
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = duration
}
