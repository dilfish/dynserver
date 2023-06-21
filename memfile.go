package main

import (
    "net/http"
    "strconv"
    "fmt"
    "time"
)


func LimitMemFileHandler(w http.ResponseWriter, r *http.Request) {
	// default value is 100 bytes
	size := uint64(100)
	u := r.RequestURI
	if len(u) > len("/limitmemfile/") {
		sizeStr := u[len("/limitmemfile/"):]
		size, _ = strconv.ParseUint(sizeStr, 10, 32)
		if size <= 0 {
			size = 100
		}
		// 10m
		if size > 10*1024*1024 {
			size = 100
		}
	}
	w.Header().Set("Content-Type", "text/html")

	str := ""
	wSize := 0
	for size >= 10 {
		wSize = wSize + 10
		size = size - 10
		// 10 bytes
		str = str + fmt.Sprintf("%09d\n", wSize)
		if len(str) > 1000 {
			w.Write([]byte(str))
			str = ""
			time.Sleep(time.Second)
		}
	}
	if size > 0 {
		for i := 1; i < int(size)+1; i++ {
			str = str + string('0'+i)
		}
	}
	if str != "" {
		w.Write([]byte(str))
	}
}

func MemFileHandler(w http.ResponseWriter, r *http.Request) {
	// default value is 100 bytes
	size := uint64(100)
	u := r.RequestURI
	if len(u) > len("/memfile/") {
		sizeStr := u[len("/memfile/"):]
		size, _ = strconv.ParseUint(sizeStr, 10, 32)
		if size <= 0 {
			size = 100
		}
		// 10m
		if size > 10*1024*1024 {
			size = 100
		}
	}
	w.Header().Set("Content-Type", "text/html")

	str := ""
	wSize := 0
	for size >= 10 {
		wSize = wSize + 10
		size = size - 10
		// 10 bytes
		str = str + fmt.Sprintf("%09d\n", wSize)
		if len(str) > 1000 {
			w.Write([]byte(str))
			str = ""
		}
	}
	if size > 0 {
		for i := 1; i < int(size)+1; i++ {
			str = str + string('0'+i)
		}
	}
	if str != "" {
		w.Write([]byte(str))
	}
}
