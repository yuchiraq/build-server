package main

import (
	"database/sql"
	"fmt"

	"build-app/base"
	"build-app/user_api"

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
