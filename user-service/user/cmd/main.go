package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/minhphuc2544/DevOps-Backend/user-service/user/internal/models"
	"github.com/minhphuc2544/DevOps-Backend/user-service/user/internal/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Hàm này lấy từng biến từ APP_ENV (inject kiểu JSON)
func getEnvFromAPP_ENV() map[string]string {
	raw := os.Getenv("APP_ENV")
	if raw == "" {
		log.Fatal("Missing APP_ENV variable")
	}
	envMap := make(map[string]string)
	err := json.Unmarshal([]byte(raw), &envMap)
	if err != nil {
		log.Fatalf("Failed to parse APP_ENV as JSON: %v\nRaw: %s", err, raw)
	}
	return envMap
}

func main() {
	env := getEnvFromAPP_ENV() // Lấy map từ APP_ENV

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		env["MYSQL_USER"],
		env["MYSQL_PASSWORD"],
		env["MYSQL_HOST"],
		env["MYSQL_PORT"],
		env["MYSQL_DATABASE"],
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
	log.Println("Successfully connected to the database.")

	router := routes.SetupRoutes(db)
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(router)))
}
