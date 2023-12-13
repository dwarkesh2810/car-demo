package request

type CreateUserRequest struct {
	FirstName string ` json:"first_name" form:"first_name" valid:"Required; MaxSize(100)"`
	LastName  string ` json:"last_name" form:"last_name" valid:"Required; MaxSize(100)"`
	Email     string ` json:"email" form:"email" valid:"Required; Email; MaxSize(100)"`
	Mobile    string ` json:"mobile" form:"mobile" valid:"Required; Mobile; MaxSize(20)"`
	Password  string ` json:"password" form:"password" valid:"Required; MaxSize(100)"`
	Role      string ` json:"role" form:"role" valid:"Required; MaxSize(20)"`
	Otp       string `json:"otp" form:"otp"`
}

type UserUpdateRequest struct {
	Id        int64  ` form:"id" json:"id"`
	FirstName string ` json:"first_name" form:"first_name" valid:"Required; MaxSize(100)"`
	LastName  string ` json:"last_name" form:"last_name" valid:"Required; MaxSize(100)"`
	Email     string ` json:"email" form:"email" valid:"Required; Email; MaxSize(100)"`
	Mobile    string ` json:"mobile" form:"mobile" valid:"Required; Mobile; MaxSize(20)"`
	Role      string ` json:"role" form:"role" valid:"Required; MaxSize(20)"`
}

type UserLoginRequest struct {
	Email    string `json:"email,omitempty" form:"email" valid:"Required; Email; MaxSize(100)"`
	Mobile   string `json:"mobile,omitempty" form:"mobile" valid:"Mobile; MaxSize(20)"`
	Password string `json:"password" form:"password" valid:"Required; MaxSize(100)"`
}

type ForgotPassword struct {
	Email       string `form:"email" json:"email" valid:"Email; MaxSize(100)"`
	Otp         string `form:"otp" json:"otp" valid:"Required"`
	NewPassword string `form:"new_password" json:"new_password" valid:"Required; MaxSize(100)"`
}

type SendOTP struct {
	Email string `form:"email" json:"email" valid:"Email; MaxSize(100)"`
}

type VerifyOTP struct {
	Email string `form:"email" json:"email" valid:"Required; Email; MaxSize(100)"`
	Otp   string `form:"otp" json:"otp" valid:"Required"`
}

type GetUserByID struct {
	Id string `json:"id"`
}

type CreateCarDataRequest struct {
	CarName string `json:"car_name" form:"car_name"`
	CarType string `json:"car_type" form:"car_type"`
	Make    string `json:"make" form:"make"`
	Model   string `json:"model" form:"model"`
}
