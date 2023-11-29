package request

type CreateUserRequest struct {
	FirstName string ` json:"first_name" form:"first_name"`
	LastName  string ` json:"last_name" form:"last_name"`
	Email     string ` json:"email" form:"email"`
	Mobile    string ` json:"mobile" form:"mobile"`
	Password  string ` json:"password" form:"password"`
	Role      string ` json:"role" form:"role"`
}

type UserRequest struct {
	Id        int64  ` form:"id" json:"id"`
	FirstName string ` json:"first_name" form:"first_name"`
	LastName  string ` json:"last_name" form:"last_name"`
	Email     string ` json:"email" form:"email"`
	Mobile    string ` json:"mobile" form:"mobile"`
	Password  string ` json:"password" form:"password"`
	Role      string ` json:"role" form:"role"`
}
