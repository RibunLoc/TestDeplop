package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/minhphuc2544/DevOps-Backend/user-service/user/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type changePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h* Handler) GetAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get all users from the database
	var users []User
	if err := h.db.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the list of users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
		"count": len(users),
		"message": "Get all users successfully",
	})
}

func (h *Handler) GetInfoByUsername(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the username from the URL parameters
	claims, err := utils.VerifyJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Logged in user ID:", claims["user_id"])

	var requestBody struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        http.Error(w, "Invalid JSON body", http.StatusBadRequest)
    }
	username := requestBody.Username
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	var user User
	if err := h.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the user as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"fullname": user.Fullname,
		"email":    user.Email,
		"created_at": user.CreatedAt,
	})
}

func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse the request body to get user details
	claims, err := utils.VerifyJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Logged in user ID:", claims["user_id"])
	
	var user changePassword
	// get the old and new password from the request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dbUser User
	// compare the old password with the hashed password in the database
	if err := h.db.Where("id = ?", claims["user_id"]).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.OldPassword)); err != nil {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	// hash the new password and update it in the database
	hashedPassword, err := utils.HashPassword(user.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if err := h.db.Model(&dbUser).Update("password", hashedPassword).Error; err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Password changed successfully",
	})
}