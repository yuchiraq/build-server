package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	// Проверка соединения с базой данных
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Connected to database!")
}

func main() {
	// Инициализация базы данных
	initDB()

	// Инициализация Gin
	r := gin.Default()

	// Обслуживание статических файлов (включая favicon)
	r.Static("/static", "./static")

	// Middleware для favicon
	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/favicon.ico" {
			c.File("./static/favicon.ico")
			c.Abort()
			return
		}
		c.Next()
	})

	// Регистрация эндпоинтов
	r.POST("/register", user_api.RegisterUser(db))
	r.GET("/register", user_api.RegisterUser(db)) // Временный GET для тестирования
	r.POST("/login", user_api.LoginUser(db))
	r.GET("/login", user_api.LoginUser(db))                    // GET для авторизации
	r.GET("/check-login", user_api.CheckLoginAvailability(db)) // Проверка логина

	// Эндпоинт для выключения сервера (GET)
	shutdown := make(chan struct{})
	r.GET("/shutdown", func(c *gin.Context) {
		log.Println("Shutdown endpoint called")
		c.JSON(http.StatusOK, gin.H{"message": "Server is shutting down..."})
		close(shutdown) // Сигнализируем о завершении работы
	})

	// Запуск сервера в отдельной горутине
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Starting server on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	select {
	case <-shutdown:
		log.Println("Shutdown signal received")
	case <-waitForInterrupt():
		log.Println("Interrupt signal received")
	}

	// Завершение работы сервера
	log.Println("Shutting down server...")
	if err := srv.Shutdown(nil); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Println("Server stopped")
}

// waitForInterrupt ожидает сигналов прерывания (Ctrl+C или SIGTERM)
func waitForInterrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	return interrupt
}
