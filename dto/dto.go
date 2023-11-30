package dto

import (
	"car_demo/models"
	"car_demo/response"
)

func DtOUserResponse(user *models.Users) *response.CreateUserResponse {
	return &response.CreateUserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Status:    user.Status,
		Role:      user.Role,
	}
}

func DtOAddCarResponse(car *models.Car_master) *response.AddCarResponse {
	return &response.AddCarResponse{
		Id:       car.Id,
		CarName:  car.CarName,
		CarImage: car.CarImage,
		Make:     car.Make,
		Model:    car.Model,
		CarType:  car.CarType,
	}
}
