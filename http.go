// http just redirect requests to https
package main

import "net/http"

type RedirectHandler struct{}

func (rd *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, 302)
}

func Redirect() {
	var rd RedirectHandler
	err := http.ListenAndServe(":80", &rd)
	if err != nil {
		panic("listen and serve http error: " + err.Error())
	}
}
