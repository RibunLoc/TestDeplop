package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/minhphuc2544/DevOps-Backend/task-service/internal/utils"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type Music struct {
	ID        uint      `json:"id"`
	Name	  string    `json:"name"`
	Artist    string   	`json:"artist"`
	// Lyrics    string   // Lyrics of the song
	AccessURL string    `json:"accessurl"` // URL to access the song
	// ImageURL  string    `gorm:"not null"` // URL to access the image of the song
	Genre	string     	 `json:"genre"`// Genre of the song
	PlayCount uint       `json:"playcount"`// Number of times the song has been played
}

type UserPlaylist struct {
	PlaylistID        uint      `json:"playlist_id"`
	UserID    uint      `json:"user_id"`
	Topic	 string    `json:"topic"`
}

type Playlist struct {
	PlaylistID        int      `json:"playlist_id"`
	MusicID    int      `json:"music_id"`
}



func (h* Handler) GetAllMusic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	claims, err := utils.VerifyJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Logged in user ID:", claims["user_id"])


	// Get all music from the database
	var music []Music
	if err := h.db.Find(&music).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the list of music as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"music": music,
		"count": len(music),
		"message": "Get all music successfully",
	})
}

func (h *Handler) UploadMusic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {


	claims, err := utils.VerifyJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Logged in user ID:", claims["user_id"])


	// get music info from request body
	var music Music
	err = json.NewDecoder(r.Body).Decode(&music)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// validate music info
	if music.Name == "" || music.Artist == "" || music.AccessURL == "" || music.Genre == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// save music to database
	err = h.db.Create(&music).Error
	if err != nil {
		http.Error(w, "Failed to save music", http.StatusInternalServerError)
		return
	}

	// return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Music uploaded successfully",
		"music":   music,
	})
	
}


func (h* Handler) CreatePlaylist (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the user ID from the request context
	claims, err := utils.VerifyJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Logged in user ID:", claims["user_id"])

	// Get the playlist information from the request body
	var playlist UserPlaylist
	err = json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the playlist information
	if playlist.UserID == 0 || playlist.Topic == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Save the playlist to the database
	err = h.db.Create(&playlist).Error
	if err != nil {
		http.Error(w, "Failed to save playlist", http.StatusInternalServerError)
		return
	}

	// Return a success response with the playlist ID
	id := playlist.PlaylistID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Playlist created successfully",
		"playlist": id,
	})
	
}

func (h* Handler) GetUserPlayListByUserId (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the user ID from the request context
	claims, err := utils.VerifyJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Logged in user ID:", claims["user_id"])

	// Get all playlists for the user from the database
	var playlists []UserPlaylist
	if err := h.db.Where("user_id = ?", claims["user_id"]).Find(&playlists).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the list of playlists as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"playlists": playlists,
	})
}

func (h* Handler) AddMusicToPlaylist (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the user ID from the request context
	claims, err := utils.VerifyJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Logged in user ID:", claims["user_id"])

	// Get the playlist ID and music ID from the request body
	var data struct {
		PlaylistID int `json:"playlist_id"`
		MusicID    int `json:"music_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the data
	if data.PlaylistID == 0 || data.MusicID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Add the music to the playlist
	err = h.db.Create(&Playlist{
		PlaylistID: data.PlaylistID,
		MusicID:    data.MusicID,
	}).Error
	if err != nil {
		http.Error(w, "Failed to add music to playlist", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Music added to playlist successfully",
	})
}

func (h *Handler) GetMusicByPlaylist(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Authenticate the user
	claims, err := utils.VerifyJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Logged in user ID:", claims["user_id"])

	// Get the playlist ID from the request parameters
	var playlist Playlist
	err = json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	playlistID := playlist.PlaylistID
	// Get the music in the playlist from the database
	var music []Music
	if err := h.db.Table("playlists").Select("*").Joins("JOIN musics ON musics.id = playlists.music_id").Where("playlists.playlist_id = ?", playlistID).Find(&music).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the list of music as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"music": music,
		"count": len(music),
		"message": "Get all music in playlist successfully",
	})

}

// incrementPlayCount increments the play count of a music item
func (h *Handler) IncrementPlayCount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Authenticate the user
	claims, err := utils.VerifyJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Logged in user ID:", claims["user_id"])

	// Get the music ID from the request body
	var PlayCount struct {
		MusicID int `json:"music_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&PlayCount)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Increment the play count in the database
	err = h.db.Model(&Music{}).Where("id = ?", PlayCount.MusicID).Update("play_count", gorm.Expr("play_count + ?", 1)).Error
	if err != nil {
		return 
	}

	var count uint
	err = h.db.Model(&Music{}).Where("id = ?", PlayCount.MusicID).Select("play_count").Scan(&count).Error
	if err != nil {
		http.Error(w, "Failed to get play count", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Play count incremented successfully",
		"music_id": PlayCount.MusicID,
		"play_count": count,
	})
}
