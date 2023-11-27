package controllers

import (
	"car_demo/conf"
	"car_demo/helper"
	"car_demo/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v5"
)

// UsersController operations for Users
type UsersController struct {
	beego.Controller
}

// URLMapping ...
func (uc *UsersController) URLMapping() {
	uc.Mapping("Post", uc.Post)
	uc.Mapping("GetOne", uc.GetOne)
	uc.Mapping("GetAll", uc.GetAll)
	uc.Mapping("Put", uc.Put)
	uc.Mapping("Delete", uc.Delete)
	uc.Mapping("Login", uc.Login)
	// uc.Mapping("VerifyUser", uc.VerifyUser)
	uc.Mapping("SendOTP", uc.SendOTP)
	uc.Mapping("VerifyOTP", uc.VerifyOTP)
	uc.Mapping("ForgetPassword", uc.ForgetPassword)
}

// Post ...
// @Title Post
// @Description create Users
// @Param	body		body 	models.Users	true		"body for Users content"
// @Success 201 {int} models.Users
// @Failure 403 body is empty
// @router / [post]
func (uc *UsersController) Post() {
	var v models.Users
	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: "+err.Error())
		return
	}

	userOTP := helper.GenerateOTP()

	up, err := helper.HashData(v.Password)

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "faild to hashed password")
		return
	}

	v.Password = up

	v.Otp = strconv.Itoa(userOTP)

	v.CreatedAt = time.Now().Unix()
	v.UpdatedAt = time.Now().Unix()

	if _, err := models.AddUsers(&v); err == nil {
		helper.JsonResponse(uc.Controller, http.StatusCreated, 1, v, "")
		helper.SendMail(v.Email, conf.EnvConfig.MailSubject, strconv.Itoa(userOTP))
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}
}

// GetOne ...
// @Title Get One
// @Description get Users by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Users
// @Failure 403 :id is empty
// @router /:id [get]
func (uc *UsersController) GetOne() {
	idStr := uc.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetUsersById(id)
	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusOK, 1, v, "")
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
// @router / [get]
func (uc *UsersController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
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
				helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error: invalid query key/value pair")
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllUsers(query, fields, sortby, order, offset, limit)
	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusOK, 1, l, "")
		return
	}
}

// Put ...
// @Title Put
// @Description update the Users
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Users	true		"body for Users content"
// @Success 200 {object} models.Users
// @Failure 403 :id is not int
// @router /:id [put]
func (uc *UsersController) Put() {
	idStr := uc.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Users{Id: id}
	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)
	if err := models.UpdateUsersById(&v); err == nil {
		helper.JsonResponse(uc.Controller, http.StatusOK, 1, "ok", "")
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}
}

// Delete ...
// @Title Delete
// @Description delete the Users
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (uc *UsersController) Delete() {
	idStr := uc.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteUsers(id); err == nil {
		helper.JsonResponse(uc.Controller, http.StatusOK, 1, "ok", "")
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}
}

func (uc *UsersController) Login() {
	var err error
	var us *models.Users
	var v models.Users

	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: "+err.Error())
		return
	}
	// json.Unmarshal(uc.Ctx.Input.RequestBody, &v)
	if v.Email == "" {
		if v.Mobile == "" {
			helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Mobile or mail requires")
			return
		} else {
			us, err = models.FindUserByMobile(v.Mobile)

			if err != nil {
				helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
				return
			}
			ok, _ := helper.VerifyHashedData(v.Password, us.Password)
			if !ok {
				helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "invalid credential")
				return
			}
		}
	} else {
		us, err = models.FindUserByEmail(v.Email)
		if err != nil {
			helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
			return
		}

		ok, _ := helper.VerifyHashedData(v.Password, us.Password)

		if !ok {
			helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "invalid credential")
			return
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": us.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	accessToken, err := token.SignedString([]byte(conf.EnvConfig.JwtSecret))

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}

	helper.JsonResponse(uc.Controller, http.StatusOK, 1, map[string]interface{}{"data": &us, "token": accessToken}, "")
}

func (uc *UsersController) ResetPassword() {
	var v struct {
		Email           string `form:"email" json:"email"`
		CurrentPassword string `form:"current_password" json:"current_password"`
		NewPassword     string `form:"new_password" json:"new_password"`
	}

	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: "+err.Error())
		return
	}

	u, err := models.FindUserByEmail(v.Email)

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "unexpected database err"+err.Error())
		return
	}

	ok, _ := helper.VerifyHashedData(v.CurrentPassword, u.Password)

	if !ok {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "email and current Password didn't match")
		return
	}

	up, err := helper.HashData(v.NewPassword)

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "faild to hashed password")
		return
	}
	u.Password = up

	// err = models.ResetPasswordUsingEmail(u)

	// if err != nil {
	// 	helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "unexpected database err"+err.Error())
	// 	return
	// }

	helper.JsonResponse(uc.Controller, http.StatusOK, 1, "OK", "")

}

func (uc *UsersController) ForgetPassword() {
	var v struct {
		Email       string `form:"email" json:"email"`
		Otp         string `form:"otp" json:"otp"`
		NewPassword string `form:"new_password" json:"new_password"`
	}

	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: "+err.Error())
		return
	}
	up, err := helper.HashData(v.NewPassword)

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "faild to hashed password")
		return
	}

	models.UserPasswordUpdate(v.Email, up, v.Otp)

	helper.JsonResponse(uc.Controller, http.StatusOK, 1, "Password updated", "")

}

func (uc *UsersController) SendOTP() {
	var v struct {
		Email string `form:"email" json:"email"`
	}

	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: "+err.Error())
		return
	}

	userOTP := helper.GenerateOTP()

	data, err := models.SendOtpToUser(v.Email, strconv.Itoa(userOTP))

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "")
		return
	}
	sent := helper.SendMail(v.Email, conf.EnvConfig.MailSubject, strconv.Itoa(userOTP))

	if !sent {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "failed to send OTP")
		return
	}

	helper.JsonResponse(uc.Controller, http.StatusOK, 1, data, "")
}

func (uc *UsersController) VerifyOTP() {
	var v struct {
		Email string `form:"email" json:"email"`
		Otp   string `form:"otp" json:"otp"`
	}

	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: ")
		return
	}
	data, err := models.VerifyUserAndUpdate(v.Email, v.Otp)

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}
	helper.JsonResponse(uc.Controller, http.StatusOK, 1, data, "")
}
