package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	campaign "github.com/syafrin-ibrahim/donasi.git/internal/app/campaign/repository"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/user/handler"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/user/repository"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/user/service"
	"github.com/syafrin-ibrahim/donasi.git/internal/package/helper"
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

	campaign := campaign.NewCampaignDBRepository(db)
	onecampaigns, err := campaign.FindByUserID(4)

	for _, newcampaign := range onecampaigns {

		fmt.Println(newcampaign.Name)
	}

	userRepo := repository.NewUserDBRepository(db)
	userService := service.NewUserService(userRepo)
	auth := middleware.NewServiceMiddleware()
	userHandler := handler.NewUserhandler(userService, auth)

	route := gin.Default()

	route.POST("register", userHandler.Register)
	route.POST("/login", userHandler.Login)
	route.POST("/email_checkers", userHandler.CheckEmailAvailability)
	route.POST("/avatars", authMiddleware(auth, userService), userHandler.UploadAvatar)

	route.Run(":8080")

}

func authMiddleware(auth middleware.Service, user handler.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			resp := helper.APIResponse("Unauthorizeed", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := auth.ValidateToken(tokenString)
		if err != nil {
			resp := helper.APIResponse("Unauthorizeed", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			resp := helper.APIResponse("Unauthorizeed", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		userID := int(claim["user_id"].(float64))
		user, err := user.GetUserByID(userID)

		if err != nil {
			resp := helper.APIResponse("Unauthorizeed", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		ctx.Set("currentUser", user)

	}
}
