package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"

	"github.com/minhphuc2544/DevOps-Backend/user-service/auth/internal/utils"
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

func (h *Handler) CreateNewUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse the request body to get user details
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the user details
	if user.Username == "" || user.Fullname == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	// Check if email is valid
	if !strings.Contains(user.Email, "@") {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	// Check if password is strong enough
	if len(user.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}
	// hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now().Format(time.RFC3339)

	// write
	if err := h.db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "username") {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "email") {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
	}

	id := user.ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Sucessfully created user",
		"id":      id,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse the request body to get user details
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the user details
	if user.Username == "" && user.Email == "" || user.Password == "" {
		http.Error(w, "Missing field", http.StatusBadRequest)
		return
	}

	var dbUser User
	if user.Username != "" {
		if err := h.db.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
	} else {
		if err := h.db.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(dbUser.ID, dbUser.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// generate a random password
	password, err := utils.GenerateRandomPassword(12)
	if err != nil {
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		return
	}
	// hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	// get the email from the request body
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the user details
	if user.Email == "" {
		http.Error(w, "Missing field", http.StatusBadRequest)
		return
	}

	// check if the email exists in the database
	var dbUser User
	if err := h.db.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}

	// update the password in the database
	if err := h.db.Model(&dbUser).Update("password", hashedPassword).Error; err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// send the password to the email
	m := gomail.NewMessage()
	m.SetHeader("From", "zacken1909@gmail.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Your New Password")
	m.SetBody("text/plain", "Your new password is: "+password)

	envPath, err := utils.LoadEnv()
	if err != nil {
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}

	err = godotenv.Load(envPath)
    if err != nil {
        log.Fatalf("Error loading .env file from %s: %v", envPath, err)
    }

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("MAIL_USERNAME") , os.Getenv("MAIL_PASSWORD"))
	err = d.DialAndSend(m)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Password reset successful. Check your email for the new password.",
	})
}