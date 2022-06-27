package main

import (
	"errors"
	"gorm.io/gorm"
	"gotest/main/models"
)

type UserControllerGetResponse struct {
	Username string `json:"username"`
}

type UserControllerPostResponse struct {
	Username string `json:"username"`
}

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		db: db,
	}
}

func (controller *UserController) Get() *UserControllerGetResponse {
	var user models.User
	controller.db.First(&user)
	return &UserControllerGetResponse{
		Username: user.Username,
	}
}

func (controller *UserController) Post() *UserControllerPostResponse {
	var user models.User
	result := controller.db.First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user := models.User{Username: "test", Password: "password"}
		result := controller.db.Create(&user)
		if result.Error != nil {
			return nil
		}

		return &UserControllerPostResponse{
			Username: user.Username,
		}
	}

	user.Username = "username"
	result = controller.db.Save(&user)
	if result.Error != nil {
		return nil
	}
	return &UserControllerPostResponse{
		Username: user.Username,
	}
}
