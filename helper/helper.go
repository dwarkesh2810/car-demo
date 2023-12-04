package helper

import (
	"car_demo/conf"
	"fmt"
	"log"
	"math/rand"

	"net/smtp"
	"os"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Success int         `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func HashData(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func VerifyHashedData(hashedString string, dataString string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(dataString), []byte(hashedString))
	check := true
	msg := ""

	if err != nil {
		msg = "password is incorrect"
		check = false
		return check, msg
	}
	return check, msg
}

func GetTokenFromHeader(c *context.Context) (string, error) {

	token := c.Input.Header("Authorization")

	// Check if the header is present and starts with "Bearer "
	if !strings.HasPrefix(token, "Bearer ") {
		return "", fmt.Errorf("invalid or Missing Token")
	}

	// Extract the token without the "Bearer " prefix
	authToken := token[7:]

	return authToken, nil

}

func GenerateOTP() int {
	rand.Seed(time.Now().UnixNano())

	// Generate a random 6-digit number
	return rand.Intn(900000) + 10000
}
func JsonResponse(c beego.Controller, statusCode int, success int, data interface{}, err string) {
	var response Response = Response{
		Success: success,
		Data:    data,
		Error:   err,
	}
	c.Ctx.Output.SetStatus(statusCode)
	c.Data["json"] = response
	c.ServeJSON()
}

func GetFileAndStore(uc beego.Controller, file string, pathName string, path string) (string, error) {
	_, hd, err := uc.GetFile(file)

	if err != nil {
		return "", err
	}

	pathForDatabase := fmt.Sprintf("%s/%s/%s/", conf.EnvConfig.BaseStoragePath, pathName, path)

	if _, err := os.Stat("./" + pathForDatabase); os.IsNotExist(err) {
		// Folder doesn't exist, so create it
		err := os.MkdirAll("./"+pathForDatabase, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	err = uc.SaveToFile(file, pathForDatabase+hd.Filename)

	if err != nil {
		return "", err
	}

	return pathForDatabase + hd.Filename, nil
}

func SendMail(to string, subject, body string) bool {
	from := conf.EnvConfig.From
	password := conf.EnvConfig.Password

	// SMTP server configuration
	smtpHost := conf.EnvConfig.SmtpHost
	smtpPort := conf.EnvConfig.SmtpPort
	// Message construction
	message := []byte("Subject: " + subject + "\r\n" + "\r\n" + body)

	// Establish a connection to the SMTP server
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

