package main

import (
	"api/auth"
	"api/campaign"
	"api/config"
	"api/handler"
	"api/helper"
	"api/payment"
	"api/transaction"
	"api/user"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	webHandler "api/web/handler"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// load the configuration
	cfg := config.InitConfig()

	//membuat koneksi ke database postgres
	dsn := cfg.DB.Host
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("connection to datatabase success")

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService(cfg)
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService(cfg, campaignRepository)
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewtransactionHandler(transactionService)

	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService, userService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	sessionWebHandler := webHandler.NewSessionHandler(userService)
	router := gin.Default()
	// Mengatur konfigurasi CORS untuk mengizinkan permintaan dari "http://localhost:3000"
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
	}))

	cookieStore := cookie.NewStore([]byte(cfg.JWTConfig.Secret))
	router.Use(sessions.Sessions("bwastartup", cookieStore))

	router.HTMLRender = loadTemplates("./web/templates")
	router.Static("/images", "./images")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterHandler)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailable)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.GET("/users", authAdminMiddleware(), userWebHandler.Index)
	router.GET("/users/new", authAdminMiddleware(), userWebHandler.New)
	router.POST("/users", authAdminMiddleware(), userWebHandler.Create)
	router.GET("/users/edit/:id", authAdminMiddleware(), userWebHandler.Edit)
	router.POST("/users/update/:id", authAdminMiddleware(), userWebHandler.Update)
	router.GET("/users/avatar/:id", authAdminMiddleware(), userWebHandler.NewAvatar)
	router.POST("/users/avatar/:id", authAdminMiddleware(), userWebHandler.CreateAvatar)

	router.GET("/campaigns", authAdminMiddleware(), campaignWebHandler.Index)
	router.GET("/campaigns/new", authAdminMiddleware(), campaignWebHandler.New)
	router.POST("/campaigns", authAdminMiddleware(), campaignWebHandler.Create)
	router.GET("/campaigns/image/:id", authAdminMiddleware(), campaignWebHandler.NewImage)
	router.POST("/campaigns/image/:id", authAdminMiddleware(), campaignWebHandler.CreateImage)
	router.GET("/campaigns/edit/:id", authAdminMiddleware(), campaignWebHandler.Edit)
	router.POST("/campaigns/update/:id", authAdminMiddleware(), campaignWebHandler.Update)
	router.GET("/campaigns/show/:id", authAdminMiddleware(), campaignWebHandler.Show)

	router.GET("/transactions", authAdminMiddleware(), transactionWebHandler.Index)

	router.GET("/login", sessionWebHandler.New)
	router.POST("/session", sessionWebHandler.Create)
	router.GET("/logout", sessionWebHandler.Destroy)

	router.Run()
}
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		var tokenString string
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userID := claims["user_id"].(float64)
		user, err := userService.GetUserByID(int(userID))
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("currentUser", user)
	}
}

func authAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userIDSession := session.Get("userID")
		if userIDSession == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
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
