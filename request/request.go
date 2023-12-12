package request

type CreateUserRequest struct {
	FirstName string ` json:"first_name" form:"first_name"`
	LastName  string ` json:"last_name" form:"last_name"`
	Email     string ` json:"email" form:"email"`
	Mobile    string ` json:"mobile" form:"mobile"`
	Password  string ` json:"password" form:"password"`
	Role      string ` json:"role" form:"role"`
}

type UserUpdateRequest struct {
	Id        int64  ` form:"id" json:"id"`
	FirstName string ` json:"first_name" form:"first_name"`
	LastName  string ` json:"last_name" form:"last_name"`
	Email     string ` json:"email" form:"email"`
	Mobile    string ` json:"mobile" form:"mobile"`
	Role      string ` json:"role" form:"role"`
}

type UserLoginRequest struct {
	Email    string `json:"email,omitempty" form:"email"`
	Mobile   string `json:"mobile,omitempty" form:"mobile"`
	Password string `json:"password" form:"password"`
}

type ForgotPassword struct {
	Email       string `form:"email" json:"email"`
	Otp         string `form:"otp" json:"otp"`
	NewPassword string `form:"new_password" json:"new_password"`
}

type SendOTP struct {
	Email string `form:"email" json:"email"`
}

type VerifyOTP struct {
	Email string `form:"email" json:"email"`
	Otp   string `form:"otp" json:"otp"`
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
