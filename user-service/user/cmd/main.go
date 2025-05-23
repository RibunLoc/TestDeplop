package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/minhphuc2544/DevOps-Backend/user-service/user/internal/models"
	"github.com/minhphuc2544/DevOps-Backend/user-service/user/internal/routes"
	"github.com/minhphuc2544/DevOps-Backend/user-service/user/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	envPath, err := utils.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load the .env file
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}
	dsn := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DATABASE")

	// Establish a database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
	log.Println("Successfully connected to the database.")

	router := routes.SetupRoutes(db) // Setup the routes
	log.Println("Starting server on :8080...")
	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(router)))
}
