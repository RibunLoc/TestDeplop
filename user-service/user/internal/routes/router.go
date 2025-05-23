package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/minhphuc2544/DevOps-Backend/user-service/user/internal/handlers"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *httprouter.Router {
	router := httprouter.New()
	h := handlers.NewHandler(db)
	// Define your routes here
	router.GET("/user/all", h.GetAllUsers)
	router.POST("/user/me", h.GetInfoByUsername)
	router.POST("/user/updatepassword", h.UpdatePassword)
	return router
}
