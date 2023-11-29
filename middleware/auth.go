package middleware

// func Middleware(c context.Context) {

// 	context.WithValue(c, "data", 1)

// 	c.Value("data")
// }

import (
	"car_demo/conf"
	"car_demo/helper"
	"car_demo/models"

	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c *context.Context) {

	tokenString, err := helper.GetTokenFromHeader(c)
	if tokenString == "" {
		c.Abort(http.StatusUnauthorized, "")
		return
	}

	if err != nil {
		c.Abort(http.StatusUnauthorized, "")
		return
	}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(conf.EnvConfig.JwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration time
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Abort(http.StatusUnauthorized, "")
			return
		}
		userId := claims["sub"].(float64)
		userIdinInt := int64(userId)

		user, err := models.GetUsersById(userIdinInt)

		if err != nil {
			c.Abort(http.StatusUnauthorized, "")
			return
		}

		if user.Id == 0 {
			c.Abort(http.StatusUnauthorized, "")
			return
		}

		c.Input.SetData("user", user)

	} else {
		c.Abort(http.StatusUnauthorized, "")
		return
	}
}
