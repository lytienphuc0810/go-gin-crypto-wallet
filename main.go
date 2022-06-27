package main

import (
	"github.com/gin-gonic/gin"
	"gotest/main/models"
	"os"
)
import "net/http"

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TestStruct struct {
	Name string `json:"name"`
}

func test() *TestStruct {
	return &TestStruct{
		Name: "name",
	}
}

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Header.Get("Accept") {
		case "application/json":
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error",
			})
		}
	}
}

func main() {
	argsWithoutProg := os.Args[1:]

	dsn := "root:password@tcp(127.0.0.1:3306)/CODEACADEMY?charset=utf8mb4&parseTime=True&loc=Local"
	if len(argsWithoutProg) >= 1 && argsWithoutProg[0] == "migrate" {
		db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		db.AutoMigrate(&models.User{})
		return
	}

	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	apiV1.Use(MiddleWare())

	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var userController = NewUserController(db)
	{
		apiV1.GET("/user", func(c *gin.Context) {
			var data = userController.Get()
			c.JSON(http.StatusOK, gin.H{
				"data": *data,
			})
		})
		apiV1.POST("/user", func(c *gin.Context) {
			var data = userController.Post()
			if data == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusBadRequest,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"data": *data,
				})
			}
		})

		apiV1.GET("/wallet", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{})
		})
		apiV1.POST("/token", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{})
		})
		apiV1.POST("/position", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{})
		})
	}

	r.Run("localhost:8080")
}
