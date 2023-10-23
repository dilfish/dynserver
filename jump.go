// 301 jump to minus 1

package main

import (
	"log"
	"net/http"
	"strconv"
)

func Jump(w http.ResponseWriter, r *http.Request) {
	if len(r.RequestURI) <= len("/jump/") {
		w.Write([]byte("jump end with 0"))
		return
	}
	jump := r.RequestURI[len("/jump/"):]
	jumpTime, err := strconv.ParseUint(jump, 10, 32)
	if err != nil {
		log.Println("bad jump time:", r.RequestURI)
		w.Write([]byte("jump end with bad jump time"))
		return
	}
	if jumpTime > 100 {
		jumpTime = 100
	}
	if jumpTime == 0 {
		w.Write([]byte("jump end with 0"))
		return
	}
	http.Redirect(w, r, "https://"+*FlagDomain+"/jump/"+strconv.FormatUint(jumpTime-1, 10), 302)
}
