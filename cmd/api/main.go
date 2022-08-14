package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/user/handler"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/user/repository"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/user/service"
	"github.com/syafrin-ibrahim/donasi.git/internal/package/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:rahasia@tcp(127.0.0.1:3306)/crowdfound?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepo := repository.NewUserDBRepository(db)
	userService := service.NewUserService(userRepo)
	auth := middleware.NewServiceMiddleware()
	userHandler := handler.NewUserhandler(userService, auth)

	route := gin.Default()

	route.POST("register", userHandler.Register)
	route.POST("/login", userHandler.Login)
	route.POST("/email_checkers", userHandler.CheckEmailAvailability)
	route.POST("/avatars", userHandler.UploadAvatar)

	route.Run(":8080")

}
