package tracks_api

import (
	"net/http"
)

func All(w http.ResponseWriter, r *http.Request) {

	/*go fmt.Println(base.TimeNow() + "||-->>" + r.RemoteAddr + " GET allTracks")
	endl := "|||"
	//code := r.URL.Query().Get("code")

	db, err := sql.Open("mysql", DataBaseConn)

	if err != nil {
		fmt.Println("Err > ", err.Error())
		return
	}

	defer db.Close()
	results, err := db.Query("SELECT id FROM tracks")

	if err != nil {
		fmt.Println("Err > ", err.Error())
		return
	}

	var tracks []base.Track
	for results.Next() {
		var trackCur base.Track
		err = results.Scan(&trackCur.ID)
		if err != nil {
			panic(err.Error())
		}
		tracks = append(tracks, trackCur)
	}

	for i := 0; i < len(tracks); i++ {
		fmt.Fprintf(w, tracks[i].ID)
		if i != len(tracks)-1 {
			fmt.Fprintf(w, endl)
		}
	}*/
}
