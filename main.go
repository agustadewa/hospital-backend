package main

import (
	"context"
	"fmt"

	"github.com/agustadewa/hospital-backend/config"
	"github.com/agustadewa/hospital-backend/controller"
	"github.com/agustadewa/hospital-backend/helper"
	"github.com/agustadewa/hospital-backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	mongoClient := helper.NewMongoConnection(ctx, config.CONFIG.MongoURL)

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", "http://localhost:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")

	r.Use(cors.New(corsConfig))

	auth := r.Group("/auth")
	{
		authController := controller.NewAuthController(mongoClient)

		auth.GET("/check", authController.CheckAuthentication)
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
	}

	admin := r.Group("/admin")
	{
		admin.Use(middleware.HeaderVerifier)

		adminController := controller.NewAdminController(mongoClient)
		admin.GET("/getaccount/:account_id", adminController.GetAccount)
		admin.POST("/createaccount", adminController.CreateAccount)

		admin.GET("/getdoctor/:doctor_id", adminController.GetDoctor)
		admin.POST("/getmanydoctor", adminController.GetManyDoctor)
		admin.POST("/createdoctor", adminController.CreateDoctor)
		admin.DELETE("/deletedoctor/:doctor_id", adminController.DeleteDoctor)
	}

	if err := r.Run(fmt.Sprintf("%s:%s", config.CONFIG.ServiceHost, config.CONFIG.ServicePort)); err != nil {
		return
	}
}
