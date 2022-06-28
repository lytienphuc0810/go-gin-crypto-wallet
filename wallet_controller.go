package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gotest/main/models"
	"log"
	"strconv"
)

type TokenResponse struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type WalletControllerGetResponse struct {
	UserID      uint            `json:"userID"`
	ID          uint            `json:"ID"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Tokens      []TokenResponse `json:"tokens"`
}

type WalletControllerAddTokenRequest struct {
	Symbol      string `json:"symbol" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type WalletControllerAddPositionRequest struct {
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
			ID:          wallet.ID,
			UserID:      wallet.UserID,
			Name:        wallet.Name,
			Description: wallet.Description,
		}
	}

	var tokens []models.Token
	controller.db.Where("wallet_id = ?", wallet.ID).First(&tokens)

	return &WalletControllerGetResponse{
		ID:          wallet.ID,
		UserID:      wallet.UserID,
		Name:        wallet.Name,
		Description: wallet.Description,
		Tokens:      controller.convertTokens(tokens),
	}
}

func (controller *WalletController) AddToken(c *gin.Context) interface{} {
	userPrincipal, _ := c.Get(identityKey)

	var user models.User
	result := controller.db.Where("username = ?", userPrincipal.(*User).UserName).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Panicln(result.Error.Error())
		return nil
	}

	walletID, _ := strconv.Atoi(c.Param("wallet_id"))
	var wallet models.Wallet
	result = controller.db.Where("id = ? AND user_id = ?", walletID, user.ID).First(&wallet)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Panicln(result.Error.Error())
		return nil
	}

	request := WalletControllerAddTokenRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Panicln(err.Error())
		return nil
	}

	token := models.Token{
		WalletID:    uint(walletID),
		Symbol:      request.Symbol,
		Name:        request.Name,
		Description: request.Description,
	}
	result = controller.db.Create(&token)
	if result.Error != nil {
		log.Panicln(result.Error.Error())
		return nil
	}
	return token
}

func (controller *WalletController) AddPosition(c *gin.Context) interface{} {
	return nil
}

func (controller *WalletController) convertTokens(tokens []models.Token) []TokenResponse {
	result := make([]TokenResponse, 0)

	for _, token := range tokens {
		result = append(result, TokenResponse{
			Symbol:      token.Symbol,
			Name:        token.Name,
			Description: token.Symbol,
		})
	}
	return result
}
