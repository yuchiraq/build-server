package main

import (
	"database/sql"
	"fmt"

	"build-app/base"
	"build_app/user_api"

	organization_api "build-app/organization_api"

	"log"
	"net/http"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

//export FreeMemory
/*func FreeMemory(pointer *int64) {
	C.free(unsafe.Pointer(pointer))
}*/

var DataBaseConn = "chiraq:ekran11Series@/building"

func main() {

	db, err := sql.Open("mysql", DataBaseConn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Инициализация Gin
	r := gin.Default()

	// Регистрация эндпоинтов
	r.POST("/register", user_api.RegisterUser(db))
	r.GET("/register", user_api.RegisterUser(db)) // Временный GET для тестирования
	r.POST("/login", user_api.LoginUser(db))
	r.GET("/check-login", user_api.CheckLoginAvailability(db)) // Проверка логина


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		go fmt.Println(base.TimeNow() + "||-->>" + r.RemoteAddr + " GET hi")
		fmt.Fprintf(w, "HI")
	})
	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "bue")
		go fmt.Println(base.TimeNow() + "||-->>" + r.RemoteAddr + " GET exit")
		os.Exit(0)
	})
	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		go fmt.Println(base.TimeNow() + "||-->>" + r.RemoteAddr + " GET notify")
		db, err := sql.Open("mysql", DataBaseConn)
		if err != nil {
			fmt.Fprintf(w, "unavailable database Tracks")
			return
		}
		defer db.Close()
		fmt.Fprintf(w, "All is good, bro")
	})
	

	// Запуск сервера
	if err := r.Run(":8090"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}
