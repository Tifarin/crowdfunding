package main

import (
	"api/handler"
	"api/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//membuat koneksi ke database postgres
	dsn := "host=localhost user=latifah password=lupakatasandi dbname=crowdfunding_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("connection to datatabase success")

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterHandler)
	api.POST("/sessions", userHandler.Login)
	router.Run()
	// //karena data banyak, disimpan di array
	// var users []user.User

	// //konek ke Db, find untuk mencari semua data di table user dan tanda "&" untuk mengambil data
	// db.Find(&users)
	// //len untuk mengetahui banyaknya data
	// length := len(users)
	// fmt.Println(length)
	// router := gin.Default()
	// router.GET("/handler", handler)
	// router.Run()
}
