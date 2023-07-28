package main

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type Msg struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Msg      string             `json:"msg" bson:"msg"`
	CreateAt time.Time          `json:"createAt" bson:"createAt"`
}

func (h *HttpsHandler) Msg(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/t/list" {
		var list []Msg
		err := h.C.Find(bson.M{}, &list)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		for _, m := range list {
			bt, _ := json.Marshal(m)
			w.Write([]byte("<h1>" + string(bt) + "</h1>"))
		}
		return
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		msg := r.Form["message"]
		message := ""
		if len(msg) > 0 {
			message = msg[0]
		}
		if len(message) > 500000 || len(message) < 2 {
			w.Write([]byte("too long or short"))
			return
		}
		var m Msg
		m.Id = primitive.NewObjectID()
		m.CreateAt = time.Now()
		m.Msg = message
		err := h.C.Insert(m)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("ok"))
		}
		return
	}
	w.Write([]byte(MessagePage))
	return
}
