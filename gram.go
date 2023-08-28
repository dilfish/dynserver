// Dilfish at Shanghai of 2020

package main

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/dilfish/telegram-bot-api-up"
	dio "github.com/dilfish/tools/io"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const Token = "1153923115:AAHUig2LQfApIF_Q-v5fn_fKgkCYhI15Flc"

func Down(u, path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Println("create file error:", path, err)
		return err
	}
	defer file.Close()
	resp, err := http.Get(u)
	if err != nil {
		log.Println("get url error:", u, err)
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println("copy error:", err)
		return err
	}
	return nil
}

func getFileName(name string) string {
	ext := filepath.Ext(name)
	if ext == "" {
		ext = ".noext"
	}
	name = dio.RandStr(5) + ext
	return name
}

func DownFile(b *tgbotapi.BotAPI, u tgbotapi.Update, path string) (string, error) {
	uri := findUrl(b, u)
	if uri == "" {
		log.Println("could not find url:", u.Message)
		return "", errors.New("could not get uri")
	}
	log.Println("find url:", uri)
	fn := getFileName(uri)
	log.Println("get fn:", fn)
	err := Down(uri, path+"/"+fn)
	if err != nil {
		log.Println("down file error:", uri, err)
		return "", err
	}
	return fn, nil
}

// findUrl got the biggest one
func findUrl(b *tgbotapi.BotAPI, u tgbotapi.Update) string {
	// not message
	if u.Message == nil {
		log.Println("empty message", "error", errors.New("empty message"))
		return ""
	}
	if len(u.Message.Photo) == 0 && u.Message.Document == nil {
		log.Println("no message photo", "error", errors.New("no photo"))
		return ""
	}
	doc := u.Message.Document
	if doc != nil {
		docu, err := b.GetFileDirectURL(doc.FileID)
		if err != nil {
			log.Println("inner error:", err)
			return ""
		}
		return docu
	}

	ps := u.Message.Photo
	// not photo
	if len(ps) == 0 {
		log.Println("not photo", "error", errors.New("not photo"))
		return ""
	}
	max := 0
	fileUrl := ""
	for _, p := range ps {
		if p.FileSize > max {
			fileUrl = p.FileID
		}
	}
	fileUrl, err := b.GetFileDirectURL(fileUrl)
	if err != nil {
		log.Println("get url error:", "error", err, "file url", fileUrl)
		return ""
	}
	return fileUrl
}

// InitTelegram create a bot for using
func InitTelegram() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Println("new bot", "error", err, "token", Token)
		return nil, nil, err
	}
	log.Println("Auth on account", "user name", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(u)
	return bot, updates, nil
}

type ReplyStruct struct {
	Now      time.Time `json:"now"`
	FileName string    `json:"fileName"`
	Error    string    `json:"error"`
}

// HandleUpdate handles two types of message
// when it is a command, we send them pictures
// when not, we return the text back
func HandleUpdate(id int64, path string, u tgbotapi.Update, bot *tgbotapi.BotAPI, fileList []string) tgbotapi.Chattable {
	log.Println("telegram message info", u.Message.From.UserName, u.Message.Text)
	var reply ReplyStruct
	reply.Now = time.Now()
	if u.Message.Document != nil || len(u.Message.Photo) > 0 {
		fn, err := DownFile(bot, u, path)
		if err != nil {
			log.Println("download file error:", err)
			reply.Error = err.Error()
		} else {
			reply.FileName = fn
		}
	} else {
		reply.Error = "does not support thie message type"
	}
	bt, _ := json.Marshal(reply)
	msg := tgbotapi.NewMessage(id, string(bt))
	return msg
}

// Telegram runs a robot
func Telegram(path string) {
	// no need to lock
	bot, updates, err := InitTelegram()
	if err != nil {
		log.Println("init telegram", "error", err)
		return
	}
	log.Println("init telegram good")

	for update := range updates {
		// ignore any non-Message Updates
		if update.Message == nil {
			log.Println("update is:", update)
			continue
		}
		msg := HandleUpdate(update.Message.Chat.ID, path, update, bot, nil)
		bot.Send(msg)
	}
}
