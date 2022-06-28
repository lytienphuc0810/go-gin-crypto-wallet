package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"index:idx_username,unique"`
	Password string
}

type Wallet struct {
	gorm.Model
	UserID      uint
	Name        string
	Description string
}

type Token struct {
	gorm.Model
	WalletID    uint   `gorm:"index:wallet_token,unique"`
	Symbol      string `gorm:"index:wallet_token,unique"`
	Name        string
	Description string
}

type Position struct {
	gorm.Model
	TokenID uint
	Amount  float64
	Note    string
}
