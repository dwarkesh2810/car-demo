package sessions

// Then initialize a global variable as the session manager:

// var globalSessions *session.Manager

// Then initialize data in your main function:

// func init() {
// 	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 3600, "providerConfig": ""}`)
// 	go globalSessions.GC()
// }
