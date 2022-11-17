package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	campaignHandler "github.com/syafrin-ibrahim/donasi.git/internal/app/campaign/handler"
	campaign "github.com/syafrin-ibrahim/donasi.git/internal/app/campaign/repository"
	campaignService "github.com/syafrin-ibrahim/donasi.git/internal/app/campaign/service"
	transactionHandler "github.com/syafrin-ibrahim/donasi.git/internal/app/transaction/handler"
	transaction "github.com/syafrin-ibrahim/donasi.git/internal/app/transaction/repository"
	transactionService "github.com/syafrin-ibrahim/donasi.git/internal/app/transaction/service"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/user/handler"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/user/repository"
	userService "github.com/syafrin-ibrahim/donasi.git/internal/app/user/service"
	webHandler "github.com/syafrin-ibrahim/donasi.git/internal/app/web/handler"
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

	//repository
	camp := campaign.NewCampaignDBRepository(db)
	new, err := camp.FindByID(1)
	fmt.Println(new)
	campaignDB := campaign.NewCampaignDBRepository(db)
	userRepo := repository.NewUserDBRepository(db)
	transactionRepo := transaction.NewTransactionDBRepository(db)

	// for _, newcampaign := range onecampaigns {
	// 	//if len(newcampaign.CampaignImages) > 0 {
	// 	fmt.Println(newcampaign.Name)
	// 	//fmt.Println(newcampaign.CampaignImages[0].FileName)

	// 	//}
	// }

	//service
	userService := userService.NewUserService(userRepo)
	camapaignService := campaignService.NewCampaignService(campaignDB)
	transactionService := transactionService.NewTransactionService(transactionRepo, campaignDB)
	auth := middleware.NewServiceMiddleware()

	// newUser := domain.Transactionparam{
	// 	ID: 1,
	// 	User: domain.User{
	// 		ID:   3,
	// 		Name: "caca",
	// 	},
	// }
	// user, err := transactionService.GetTransactionByCampaignID(newUser)
	// fmt.Println(user)
	// campaign := domain.CreateCampaignParam{}
	// campaign.Name = "Dana Talangan"
	// campaign.ShortDescription = "short talangan"
	// campaign.Description = "long talangan"
	// campaign.GoalAmount = 1000000000
	// campaign.Perks = "satu, dua, tiga"
	// newUser, _ := userService.GetUserByID(3)
	// campaign.User = newUser

	// _, err = camapaignService.CreateCampaign(campaign)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	//handler
	userHandler := handler.NewUserhandler(userService, auth)
	campaignHandler := campaignHandler.NewCampaignHandler(camapaignService)
	transactionHandler := transactionHandler.NewTransactionHandler(transactionService)

	//web
	userWebHandler := webHandler.NewUserHandler(userService)

	route := gin.Default()
	//route.LoadHTMLGlob("../../internal/app/web/templates/**/*")
	route.Static("/images", "./internal/app/images")
	route.HTMLRender = loadTemplates("internal/app/web/templates")
	route.Static("/css", "./internal/app/web/assets/css")
	route.Static("/js", "./internal/app/web/assets/js")
	route.Static("/webfonts", "./internal/app/web/assets/webfonts")

	route.POST("register", userHandler.Register)
	route.POST("/login", userHandler.Login)
	route.POST("/email_checkers", userHandler.CheckEmailAvailability)
	route.POST("/avatars", authMiddleware(auth, userService), userHandler.UploadAvatar)
	route.GET("/campaigns", campaignHandler.GetCampaigns)
	route.POST("/campaign", authMiddleware(auth, userService), campaignHandler.CreateCampaign)
	route.GET("/campaign/:id", campaignHandler.GetCampaign)
	route.PUT("/campaign/:id", authMiddleware(auth, userService), campaignHandler.UpdateCampaign)
	route.POST("/campaign-images", authMiddleware(auth, userService), campaignHandler.UploadImage)
	route.GET("/campaigns/:id/transactions", authMiddleware(auth, userService), transactionHandler.GetCampaignTransaction)
	route.GET("/transactions", authMiddleware(auth, userService), transactionHandler.GetUserTransactions)

	//web
	route.GET("/users/", userWebHandler.Index)
	route.GET("/users/new", userWebHandler.New)
	route.POST("/users", userWebHandler.Create)
	route.GET("/users/edit/:id", userWebHandler.Edit)
	route.POST("/users/update/:id", userWebHandler.Update)

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

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
