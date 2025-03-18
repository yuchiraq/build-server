package main

import (
	"database/sql"
	"fmt"
	"log"

	"build-app/user_api"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// initDB инициализирует подключение к базе данных
func initDB() {
	var err error
	db, err = sql.Open("mysql", "chiraq:ekran11Series@/building")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Connected to database!")
}

// setupRoutes настраивает маршруты для Gin
func setupRoutes(r *gin.Engine) {
	// Обслуживание статических файлов
	r.Static("/static", "./static")

	// Обработка favicon
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./static/favicon.ico")
	})

	// Эндпоинты для работы с пользователями
	r.POST("/register", user_api.RegisterUser(db))
	r.GET("/register", user_api.RegisterUser(db)) // Временный GET для тестирования
	r.POST("/login", user_api.LoginUser(db))
	r.GET("/login", user_api.LoginUser(db))
	r.GET("/check-login", user_api.CheckLoginAvailability(db))
}

func main() {
	// Инициализация базы данных
	initDB()

	// Инициализация Gin
	r := gin.Default()

	// Настройка маршрутов
	setupRoutes(r)

	// Запуск сервера на порту 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
