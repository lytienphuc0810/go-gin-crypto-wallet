package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gotest/main/models"
	"log"
	"net/http"
	"os"
)

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Header.Get("Accept") {
		case "application/json":
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Only Accepting JSON",
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
		db.AutoMigrate(&models.Wallet{})
		db.AutoMigrate(&models.Token{})
		db.AutoMigrate(&models.Position{})
		db.Unscoped().Where("deleted_at IS NULL").Delete(&models.User{})
		db.Create(&models.User{
			Username: "test1",
		})
		db.Create(&models.User{
			Username: "test2",
		})
		db.Create(&models.User{
			Username: "test3",
		})
		db.Create(&models.User{
			Username: "test4",
		})
		db.Create(&models.User{
			Username: "test5",
		})
		return
	}

	r := gin.Default()
	authMidlw, err := initJWTAuth(r)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	apiV1 := r.Group("/api/v1")
	apiV1.Use(MiddleWare(), authMidlw)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var userController = NewUserController(db)
	var walletController = NewWalletController(db)
	{
		apiV1.GET("/profile", func(c *gin.Context) {
			var profile = userController.Get(c)
			c.JSON(http.StatusOK, gin.H{"profile": profile})
		})
		apiV1.GET("/wallet", func(c *gin.Context) {
			var wallet = walletController.Get(c)
			c.JSON(http.StatusOK, gin.H{"wallet": wallet})
		})
		apiV1.POST("/token", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{})
		})
		apiV1.POST("/position", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{})
		})
	}
	err = r.Run("localhost:8080")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

/*
 * /login
 * /auth/refresh_token
 */
func initJWTAuth(r *gin.Engine) (gin.HandlerFunc, error) {
	authMiddleware, err := GetAuthMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
		return nil, err
	}
	err = authMiddleware.MiddlewareInit()

	if err != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
		return nil, err
	}

	r.POST("/login", authMiddleware.LoginHandler)

	//r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
	//	claims := jwt.ExtractClaims(c)
	//	log.Printf("NoRoute claims: %#v\n", claims)
	//	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	//})

	auth := r.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	return authMiddleware.MiddlewareFunc(), nil
}
