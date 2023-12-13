package controllers

import (
	"car_demo/conf"
	"car_demo/dto"
	"car_demo/helper"
	"car_demo/logger"
	"car_demo/models"
	"car_demo/request"

	"encoding/json"
	"log"

	"net/http"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/golang-jwt/jwt/v5"
)

// UsersController operations for Users
type UsersController struct {
	beego.Controller
	i18n.Locale
}

// Register ...
// @Title Post
// @Description create Users
// @Param	body		body 	request.CreateUserRequest	true		"body for Users content"
// @Success 201 {int} models.Users
// @Failure 403 body is empty
// @router /users/register [post]
func (uc *UsersController) Register() {
	logger.Init()
	var v request.CreateUserRequest
	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.parsing"))
		return
	}

	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)

	isExist := models.IsExistingUser(v.Email, v.Mobile)

	if isExist {
		logger.Warning(helper.LanguageTranslate(uc.Controller, "error.userexist"), v.Email)
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.userexist"))
		return
	}
	userOTP := helper.GenerateOTP()
	v.Otp = string(rune(userOTP))

	if _, err := models.AddUsers(&v); err == nil {
		go helper.SendMail(v.Email, conf.EnvConfig.MailSubject, strconv.Itoa(userOTP))
		T, _ := models.LastInsertedUser()
		helper.JsonResponse(uc.Controller, http.StatusCreated, 1, helper.MappedData("user is created", dto.DtOUserResponse(T)), "")
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.db"))
		return
	}
}

// GetOne ...
// @Title Get One
// @Description get Users by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Users
// @Failure 403 :id is empty
// @router /users/getone/:id([0-9]+) [get]
func (uc *UsersController) GetOne() {

	i := uc.Ctx.Input.Params()

	id, _ := strconv.ParseInt(i["0"], 0, 64)
	// id := int64(220)
	v, err := models.GetUsersById(id)
	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.db"))
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusOK, 1, helper.MappedData("user's data", v), "")
		return
	}
}

// GetAll ...
// @Title Get All
// @Description get Users
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Users
// @Failure 403
// @Failure 400
// @router /users/getall [get]
func (uc *UsersController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 1
	var offset int64

	// fields: col1,col2,entity.col3
	if v := uc.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := uc.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := uc.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := uc.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := uc.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := uc.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.query"))
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllUsers(query, fields, sortby, order, offset, limit)
	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.db"))
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusOK, 1, helper.MappedData("all users data", l), "")
		return
	}
}

// Put ...
// @Title Put
// @Description update the Users
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	request.UserUpdateRequest	true		"body for Users content"
// @Success 200 {object} string
// @Failure 403 :id is not int
// @Failure 400
// @router /users/update [put]
func (uc *UsersController) Update() {

	var v request.UserUpdateRequest
	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.parsing"))
		return
	}

	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)

	if err := models.UpdateUsersById(&v); err == nil {
		helper.JsonResponse(uc.Controller, http.StatusOK, 1, helper.MappedData("user's data is updated", nil), "")
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.db"))
		return
	}
}

// Deletes ...
// @Title Delete
// @Description delete the Users
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @Failure 400
// @router /users/deletes/:id([0-9]+) [delete]
func (uc *UsersController) Deletes() {
	idStr := uc.Ctx.Input.Params()
	id, _ := strconv.ParseInt(idStr["0"], 0, 64)

	if err := models.DeleteUsers(id); err == nil {
		helper.JsonResponse(uc.Controller, http.StatusOK, 1, helper.MappedData("user's data is deleted", nil), "")
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.db"))
		return
	}
}

// Login ...
// @Title Post
// @Description Login User
// @Param	body		body 	request.UserLoginRequest	true		"body for Users content"
// @Success 200 {int} map[string]interface{}
// @Failure 403 body is empty
// @Failure 400
// @router /users/login [post]
func (uc *UsersController) Login() {
	var err error
	var us *models.Users
	var v request.UserLoginRequest

	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.parsing"))
		return
	}

	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)
	if v.Email == "" {
		if v.Mobile == "" {
			helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.login"))
			return
		} else {
			us, err = models.FindUserByMobile(v.Mobile)

			if err != nil {
				helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.db"))
				return
			}
			ok, _ := helper.VerifyHashedData(v.Password, us.Password)
			if !ok {
				helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.credential"))
				return
			}
		}
	} else {
		us, err = models.FindUserByEmail(v.Email)
		if err != nil {
			helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.db"))
			return
		}

		ok, _ := helper.VerifyHashedData(v.Password, us.Password)

		if !ok {
			helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.credential"))
			return
		}
	}
	// err = sessions.Set(uc.Controller, "user", us.Email)

	// if err != nil {
	// 	helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.setsession"))
	// 	return
	// }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": us.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	accessToken, err := token.SignedString([]byte(conf.EnvConfig.JwtSecret))

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.accesstoken"))
		return
	}

	helper.JsonResponse(uc.Controller, http.StatusOK, 1, helper.MappedData("user is logged in", accessToken), "")
}

// ForgetPassword ...
// @Title Post
// @Description Forgot Password
// @Param	body		body 	request.ForgotPassword	true		"body for Users content"
// @Success 200 {int} string
// @Failure 403 body is empty
// @Failure 400
// @router /users/forgot_password [post]
func (uc *UsersController) ForgetPassword() {
	var v request.ForgotPassword

	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.parsing"))
		return
	}
	up, err := helper.HashData(v.NewPassword)

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.hashpassword"))
		return
	}

	models.UserPasswordUpdate(v.Email, up, v.Otp)

	helper.JsonResponse(uc.Controller, http.StatusOK, 1, "Password updated", "")

}

// SendOTP ...
// @Title SendOTP
// @Description Send Otp
// @Param	body		body 	request.SendOTP	true		"body for Users content"
// @Success 200 {int} string
// @Failure 403 body is empty
// @Failure 400
// @router /users/sendotp [post]
func (uc *UsersController) SendOTP() {
	var v request.SendOTP

	if err := uc.ParseForm(&v); err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.parsing"))
		return
	}
	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)

	userOTP := helper.GenerateOTP()

	_, err := models.SendOtpToUser(v.Email, strconv.Itoa(userOTP))

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "")
		return
	}
	sent, _ := helper.SendMail(v.Email, conf.EnvConfig.MailSubject, strconv.Itoa(userOTP))

	if !sent {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.sendotp"))
		return
	}

	helper.JsonResponse(uc.Controller, http.StatusOK, 1, "otp is sent", "")
}

// VerifyOTP ...
// @Title Post
// @Description Verify OTP
// @Param	body		body 	request.VerifyOTP	true		"body for Users content"
// @Success 200 {int} string
// @Failure 403 body is empty
// @Failure 400
// @router /users/verifyotp [post]
func (uc *UsersController) VerifyOTP() {
	var v request.VerifyOTP

	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.parsing"))
		return
	}
	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)
	_, err := models.VerifyUserAndUpdate(v.Email, v.Otp)

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, helper.LanguageTranslate(uc.Controller, "error.db"))
		return
	}
	helper.JsonResponse(uc.Controller, http.StatusOK, 1, "user's verification is completed", "")
}

// DemoSet ...
// @Title DemoSet
// @Description Demoset
// @Success 200 {int} string
// @Failure 403 body is empty
// @Failure 400
// @router /users/demoset [post]
func (c *UsersController) DemoSet() {
	err := c.SetSession("test", "hello")
	if err != nil {
		helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, err.Error())
		log.Print(err)
		return
	}
	err = c.SetSession("test2", "how ar you")

	c.StartSession()

	if err != nil {
		helper.JsonResponse(c.Controller, http.StatusBadRequest, 0, nil, err.Error())
		log.Print(err)
	}

	helper.JsonResponse(c.Controller, http.StatusOK, 1, "set data", "")
}

// DemoGet ...
// @Title DemoGet
// @DescriptionDemoGet
// @Param   Accept-Language  header  string  false  "Bearer YourAccessToken"
// @Success 200 {int} string
// @Failure 403 body is empty
// @Failure 400
// @router /users/demoget [get]
func (c *UsersController) DemoGet() {
	adata := helper.LanguageTranslate(c.Controller, "error.db")
	bdata := helper.LanguageTranslate(c.Controller, "bye.bye")
	helper.JsonResponse(c.Controller, 200, 1, map[string]interface{}{"a": adata, "b": bdata}, "")
}
