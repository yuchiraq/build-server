package main

import (
	"database/sql"
	"fmt"

	"build-app/base"
	_ "build-app/base"

	organization_api "build-app/organization_api"

	//_ "github.com/bobertlo/go-mpg123/mpg123"

	"log"
	"net/http"
	"os"
	//_ "github.com/go-sql-driver/mysql"
)

//export FreeMemory
/*func FreeMemory(pointer *int64) {
	C.free(unsafe.Pointer(pointer))
}*/

// DataBaseConn /*tcp(34.69.28.110)*/
var DataBaseConn = "root:ekran11series@/beats"

func main() {

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
	http.HandleFunc("/tracks/all", organization_api.All)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
