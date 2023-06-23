// Sean at shanghai
// 2020
// print cloudflare info

package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

type IPInfo struct {
	Key   string
	Value string
}

// CFIPHandler handles cloudflare HTTP headers info
func CFIPHandler(w http.ResponseWriter, r *http.Request) {
	var headerList []string
	var kvList []IPInfo
	for k, _ := range w.Header() {
		headerList = append(headerList, k)
	}
	sort.Strings(headerList)
	for _, hdr := range headerList {
		var ii IPInfo
		vv := strings.Join(w.Header()[hdr], ",")
		ii.Key = hdr
		ii.Value = vv
		kvList = append(kvList, ii)
	}
	if r.Method != "GET" {
		bt, _ := json.Marshal(kvList)
		w.Write(bt)
		return
	}
	RenderHTML("ip.index.html", w, kvList)
	return
}

func RenderHTML(fn string, wr io.Writer, data interface{}) error {
	t, err := template.ParseFiles(fn)
	if err != nil {
		log.Println("parse file", "error", err, "file_name", fn)
		return err
	}
	return t.Execute(wr, data)
}
