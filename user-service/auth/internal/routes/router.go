package routes

import (

	"github.com/julienschmidt/httprouter"
	"github.com/minhphuc2544/DevOps-Backend/user-service/auth/internal/handlers"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *httprouter.Router {
	router := httprouter.New()
	h := handlers.NewHandler(db)
	// Define your routes here
	router.POST("/auth/signup", h.CreateNewUser)
	router.POST("/auth/login", h.Login)
	router.POST("/auth/forgotpassword", h.ForgotPassword)
	return router
}