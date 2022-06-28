package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gotest/main/models"
)

type WalletControllerGetResponse struct {
	UserID      uint   `json:"userID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type WalletController struct {
	db *gorm.DB
}

func NewWalletController(db *gorm.DB) *WalletController {
	return &WalletController{
		db: db,
	}
}

func (controller *WalletController) Get(c *gin.Context) *WalletControllerGetResponse {
	userPrincipal, _ := c.Get(identityKey)

	var user models.User
	result := controller.db.Where("username = ?", userPrincipal.(*User).UserName).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	var wallet models.Wallet
	result = controller.db.Where("user_id = ?", user.ID).First(&wallet)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		wallet := models.Wallet{
			UserID:      user.ID,
			Name:        "Wallet",
			Description: "Sample Wallet",
		}
		result := controller.db.Create(&wallet)
		if result.Error != nil {
			return nil
		}
		return &WalletControllerGetResponse{
			UserID:      wallet.UserID,
			Name:        wallet.Name,
			Description: wallet.Description,
		}
	}
	return &WalletControllerGetResponse{
		UserID:      wallet.UserID,
		Name:        wallet.Name,
		Description: wallet.Description,
	}
}
