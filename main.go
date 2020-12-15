package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TOKEN  string
	CHATID string
	PORT   string `default:"6666"`
}

var bot *tgbotapi.BotAPI
var port string
var curchatid int64

func send_msg_to_current(text string) {
	msg := tgbotapi.NewMessage(curchatid, text)
	if _, err := bot.Send(msg); err != nil {
		log.Fatal(err.Error())
	}
}

func send_log(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	logpath := query.Get("path")

	msg := tgbotapi.NewDocumentUpload(curchatid, logpath)
	if _, err := bot.Send(msg); err != nil {
		log.Fatal(err.Error())
	}
}

func send_txt(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	text := query.Get("text")

	if text != "" {
		send_msg_to_current(text)
	}
}

func main() {
	var c Config

	botname := "logbot"
	if len(os.Args) > 1 {
		botname = string(os.Args[1])
	}

	err := envconfig.Process(botname, &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	port = c.PORT
	curchatid, _ = strconv.ParseInt(c.CHATID, 10, 64)

	bot, err = tgbotapi.NewBotAPI(c.TOKEN)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Bot name: " + bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err.Error())
	}

	http.HandleFunc("/send_txt", send_txt)
	http.HandleFunc("/send_log", send_log)
	go func() {
		log.Fatal(http.ListenAndServe("127.0.0.1:"+c.PORT, nil))
	}()

	for update := range updates {
		log.Printf("Got msg: %s from chat %d\n", update.Message.Text, update.Message.Chat.ID)
	}
}
