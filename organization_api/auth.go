package organization_api 

import (
	"database/sql"
	"net/http"
	"fmt"

	"build-app/base"

	_ "github.com/go-sql-driver/mysql"
)

func Auth(w http.ResponseWriter, r *http.Request) {

	go fmt.Println(base.TimeNow() + "||-->>" + r.RemoteAddr + " GET auth")
	endl := "|||"
	login := r.URL.Query().Get("login")
	//password := r.URL.Query().Get("pass")

	db, err := sql.Open("mysql", DataBaseConn)

	if err != nil {
		fmt.Println("Err > ", err.Error())
		return
	}

	defer db.Close()
	results, err := db.Query("SELECT password FROM users WHERE login = '" + login + "'")

	if err != nil {
		fmt.Println("Err > ", err.Error())
		return
	}

	var users []base.User
	for results.Next() {
		var strCur base.User
		err = results.Scan(&strCur.Password)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, strCur)
	}

	for i := 0; i < len(users); i++ {
		fmt.Fprintf(w, users[i].Password)
		if i != len(users)-1 {
			fmt.Fprintf(w, endl)
		}
	}
}
