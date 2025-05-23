package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/minhphuc2544/DevOps-Backend/task-service/internal/handlers"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *httprouter.Router {
	router := httprouter.New()
	h := handlers.NewHandler(db)
	// Define your routes here
	router.GET("/task/getAllMusic", h.GetAllMusic) // Get all music
	router.POST("/task/uploadMusic", h.UploadMusic) // Upload music
	router.POST("/task/createPlaylist", h.CreatePlaylist) // Create playlist
	router.POST("/task/addMusicToPlaylist", h.AddMusicToPlaylist) // Add music to playlist
	router.GET("/task/getUserPlaylist", h.GetUserPlayListByUserId) // Get playlist
	router.GET("/task/getMusicInPlaylist", h.GetMusicByPlaylist) // Get music by playlist id
	router.POST("/task/incrementPlayCount", h.IncrementPlayCount) // Increment play count
	return router
}
