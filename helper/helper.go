package helper

import (
	"car_demo/conf"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"reflect"

	"net/smtp"
	"os"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/i18n"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Success int         `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Now(duration int64) string {
	t := time.UnixMilli(duration)
	return t.Format("2006-01-02 15:04:05")
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
	return rand.Intn(900000) + 10000
}
func JsonResponse(c beego.Controller, statusCode int, success int, data interface{}, err interface{}) {
	var response Response = Response{
		Success: success,
		Data:    data,
		Error:   err,
	}
	c.Ctx.Output.SetStatus(statusCode)
	c.Data["json"] = response
	c.ServeJSON()
}

func GetFileExtensionFromForm(c *context.Context, file string) (string, error) {
	_, hd, err := c.Request.FormFile(file)
	if err != nil {
		return "", err
	}

	fileName := hd.Filename

	splitFileName := strings.Split(fileName, ".")

	return splitFileName[1], nil
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

func SendMail(to string, subject, body string) (bool, error) {
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
		return false, err
	}
	return true, nil
}

func MappedData(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{"message": message, "data": data}
}

func LanguageTranslate(c beego.Controller, key string) string {
	lang := c.Ctx.Input.GetData("lang").(string)
	language := strings.ToLower(lang)
	switch language {
	case "en-us", "en", "en-U", "US":
		lang = "en-US"
	case "hi-in", "hi", "hi-I", "IN":
		lang = "hi-IN"
	case "zh-cn", "zh", "zh-C", "CN":
		lang = "zh-CN"
	}
	return i18n.Tr(lang, key)
}

func SecondsToMinutesAndSeconds(seconds int64) (int64, int64) {
	minutes := seconds / 60
	remainingSeconds := seconds % 60
	return minutes, remainingSeconds
}

func SecondsToDayHourMinAndSeconds(seconds int) (int64, int64, int64, int64) {
	days := seconds / 86400
	hour := (seconds % 86400) / 3600
	minute := (seconds % 3600) / 60
	second := seconds % 60
	return int64(days), int64(hour), int64(minute), int64(second)
}

func GenerateModel(field interface{}, name string) bool {
	modelType := reflect.TypeOf(field)

	// Construct the -fields string based on the struct fields and types
	var fields []string
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fields = append(fields, fmt.Sprintf("%s:%s", field.Name, field.Type))
	}

	// Join the fields into a comma-separated string
	fieldsStr := strings.Join(fields, ",")

	// Run the bee generate model command with dynamically generated -fields argument
	cmd := exec.Command("bee", "generate", "model", name, "-fields="+fieldsStr)

	// Set output capture
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Println("Exit status:", exitErr.ExitCode())
			fmt.Println("Command output:", string(cmdOutput))
			return false
		}
		return false
	}
	return true
}

func GenerateController(name string) bool {

	// Run the bee generate model command with dynamically generated -fields argument
	cmd := exec.Command("bee", "generate", "controller", name)

	// Set output capture
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Println("Exit status:", exitErr.ExitCode())
			fmt.Println("Command output:", string(cmdOutput))
			return false
		}
		return false
	}
	return true
}

func GenerateMigration(field interface{}, name string, driver string, conn string) bool {
	modelType := reflect.TypeOf(field)

	// Construct the -fields string based on the struct fields and types
	var fields []string
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fields = append(fields, fmt.Sprintf("%s:%s", field.Name, field.Type))
	}

	// Join the fields into a comma-separated string
	fieldsStr := strings.Join(fields, ",")

	// Run the bee generate model command with dynamically generated -fields argument
	cmd := exec.Command("bee", "generate", "migration", name, "-fields="+fieldsStr, "-driver="+driver, "-conn="+conn)

	fmt.Println(cmd)

	// Set output capture
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Println("Exit status:", exitErr.ExitCode())
			fmt.Println("Command output:", string(cmdOutput))
			return false
		}
		return false
	}
	return true
}

func SliceToString(data []string) string {
	return strings.Join(data, ",")
}
