package response

import "car_demo/models"

type CreateUserResponse struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	Status    int    `json:"status"`
	Role      string `json:"role"`
}

type AddCarResponse struct {
	Id       int64          `json:"id"`
	UserId   int64          `json:"user_id"`
	CarName  string         `json:"car_name"`
	CarImage string         `json:"car_image" `
	Make     string         `json:"make"`
	Model    string         `json:"model"`
	CarType  models.CarType `json:"car_type"`
}
