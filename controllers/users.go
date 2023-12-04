package controllers

import (
	"car_demo/conf"
	"car_demo/dto"
	"car_demo/helper"
	"car_demo/models"
	"car_demo/request"
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
// func (uc *UsersController) URLMapping() {
// 	uc.Mapping("Post", uc.Post)
// 	uc.Mapping("GetOne", uc.GetOne)
// 	uc.Mapping("GetAll", uc.GetAll)
// 	uc.Mapping("Put", uc.Put)
// 	uc.Mapping("Delete", uc.Delete)
// 	uc.Mapping("Login", uc.Login)
// 	uc.Mapping("SendOTP", uc.SendOTP)
// 	uc.Mapping("VerifyOTP", uc.VerifyOTP)
// 	uc.Mapping("ForgetPassword", uc.ForgetPassword)
// }

// Register ...
// @Title Post
// @Description create Users
// @Param	body		body 	request.CreateUserRequest	true		"body for Users content"
// @Success 201 {int} models.Users
// @Failure 403 body is empty
// @router /users/register [post]
func (uc *UsersController) Register() {
	var v request.CreateUserRequest
	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: "+err.Error())
		return
	}

	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)

	isExist := models.IsExistingUser(v.Email, v.Mobile)
	if isExist {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "User already exist")
		return
	}

	if _, err := models.AddUsers(&v); err == nil {
		// 	helper.JsonResponse(uc.Controller, http.StatusCreated, 1, u, "")
		// helper.SendMail(v.Email, conf.EnvConfig.MailSubject, strconv.Itoa(userOTP))

		T, _ := models.LastInsertedUser()
		helper.JsonResponse(uc.Controller, http.StatusCreated, 1, dto.DtOUserResponse(T), "")
		return
	} else {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}
}

// // GetOne ...
// // @Title Get One
// // @Description get Users by id
// // @Param	body		body 	request.GetUserByID	true		"body for Users content"
// // @Success 200 {object} response.CreateUserResponse
// // @Failure 403  is empty
// // @router /users/getone [post]
// func (uc *UsersController) GetOne() {
// 	var v request.GetUserByID

// 	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)

// 	id, _ := strconv.ParseInt(v.Id, 0, 64)
// 	u, err := models.GetUsersById(id)
// 	if err != nil {
// 		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
// 		return
// 	} else {
// 		helper.JsonResponse(uc.Controller, http.StatusOK, 1, dto.DtOUserResponse(u), "")
// 		return
// 	}
// }

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
// @Param	body		body 	request.UserUpdateRequest	true		"body for Users content"
// @Success 200 {object} string
// @Failure 403 :id is not int
// @Failure 400
// @router /users/:id [put]
func (uc *UsersController) Put() {
	var v request.UserUpdateRequest
	if err := uc.ParseForm(&v); err != nil {
		// Handle error if parsing fails
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: "+err.Error())
		return
	}

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
// @Failure 400
// @router /users/:id [delete]
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
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: "+err.Error())
		return
	}

	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)
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

	// if err := uc.ParseForm(&v); err != nil {
	// 	helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: "+err.Error())
	// 	return
	// }
	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)

	userOTP := helper.GenerateOTP()

	_, err := models.SendOtpToUser(v.Email, strconv.Itoa(userOTP))

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "")
		return
	}
	sent := helper.SendMail(v.Email, conf.EnvConfig.MailSubject, strconv.Itoa(userOTP))

	if !sent {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "failed to send OTP")
		return
	}

	helper.JsonResponse(uc.Controller, http.StatusOK, 1, "Otp Sent Successfully", "")
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
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, "Error while parsing form data: ")
		return
	}
	json.Unmarshal(uc.Ctx.Input.RequestBody, &v)
	_, err := models.VerifyUserAndUpdate(v.Email, v.Otp)

	if err != nil {
		helper.JsonResponse(uc.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}
	helper.JsonResponse(uc.Controller, http.StatusOK, 1, "User Verification success", "")
}
