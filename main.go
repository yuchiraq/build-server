package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	//"build-app/base"
	"build-app/user_api"

	//organization_api "build-app/organization_api"

	"log"

	"github.com/gin-gonic/gin"

	//"os"
	_ "github.com/go-sql-driver/mysql"
)

//export FreeMemory
/*func FreeMemory(pointer *int64) {
	C.free(unsafe.Pointer(pointer))
}*/

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "chiraq:ekran11Series@/building")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database!")
}

func main() {

	initDB()

	// Инициализация Gin
	r := gin.Default()

	// Регистрация эндпоинтов
	r.POST("/register", user_api.RegisterUser(db))
	r.GET("/register", user_api.RegisterUser(db)) // Временный GET для тестирования
	r.POST("/login", user_api.LoginUser(db))
	r.GET("/login", user_api.LoginUser(db))                    // GET для авторизации
	r.GET("/check-login", user_api.CheckLoginAvailability(db)) // Проверка логина

	// Эндпоинт для выключения сервера (GET)
	shutdown := make(chan struct{})
	r.GET("/shutdown", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server is shutting down..."})
		close(shutdown) // Сигнализируем о завершении работы
	})

	// Запуск сервера в отдельной горутине
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
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
