package main

import (
	// "crypto/rand"
	"math/rand"
	"database/sql"
	"fmt"
	// "io"
	"log"
	// "strconv"
	"strings"
	// "time"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	db, err := sql.Open(dbDriver, "b0c87dd85bc29f:8c07d270ec2548f@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_d707738d1e2741f")
	if err != nil {
		fmt.Println("KOOOOOOOOOL")
		panic(err.Error())
	}
	return db
}

// var buttons = tgbotapi.NewInlineKeyboardMarkup(
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonText("Личность"),
// 		tgbotapi.NewInlineKeyboardButtonText("Тело"),
// 		tgbotapi.NewInlineKeyboardButtonText("Отношения"),
// 	),
// )

var buttons = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Личность",),
		tgbotapi.NewKeyboardButton("Тело"),
		tgbotapi.NewKeyboardButton("Отношения"),
	),
)


type Message struct {
	Message string
	Category string
}

var ctg string

func main() {

	//create tables
	db := dbConn()
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS usenmesgs (id INTEGER AUTO_INCREMENT PRIMARY KEY, message TEXT, category TEXT)")
	if err != nil {
		fmt.Println(err)
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS lastmesgs (id INTEGER AUTO_INCREMENT PRIMARY KEY, message TEXT)")
	if err != nil {
		fmt.Println(err)
	}
	statement.Exec()
	db.Close()


	//connect to bot
	bot, err := tgbotapi.NewBotAPI("1379112080:AAF0CHbABaAmJQL0xcrvhHzf5OVU_4eBgjs")
	if err != nil {
		log.Panic(err)
	}

	//collect  and send updates
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue 
		}

		if update.CallbackQuery != nil {
			fmt.Println("8888")
			fmt.Println("THIS IS II: ", update.CallbackQuery.Data)
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}

		if update.Message.Text == "Личность" || update.Message.Text == "Тело" || update.Message.Text == "Отношения"{
			db := dbConn()
			// message := strings.Split(update.Message.Text, "/usen")
			// message = strings.Split(message[1], " ")
			// category := message[1]
			var ms []string
			rows, err := db.Query("SELECT message FROM usenmesgs WHERE category = ?", update.Message.Text)
			if err == nil {
				for rows.Next() {
					var m string
					if err := rows.Scan(&m); err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command :(")
						bot.Send(msg)
						break
					}
					ms = append(ms, m)
				}
			}
			
			//select random message from category
			if len(ms) > 0 {
				index := rand.Intn(len(ms))
				fmt.Println(ms[index])

			
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, ms[index])
				bot.Send(msg)

		
				
			}
		} else if update.Message.IsCommand() == false {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не прочь поболтать, но я простой бот, и пока не умею так!\nНо @usorohpaius всегда рад сообщениям!")
			bot.Send(msg)
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				username := update.Message.From.UserName
				if username == "dinadinus" || username == "usorohpaius" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, Дина!\nПро что ты хочешь комплимент?")
				msg.ReplyMarkup = buttons
				bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Прости! Этот бот только для Динуса.\nА тебе хорошего дня!")
					bot.Send(msg)
				}
			case "insert":
				db := dbConn()
				message := strings.Split(update.Message.Text, "/insert")
				message = strings.Split(message[1], " ")
				category := message[1]
				text := strings.Join(message[2:], " ")
				fmt.Println("text:", text)
				fmt.Println("category:", category)
				fmt.Println(len(text))

				if len(text) > 0 {
					statement, err := db.Prepare("INSERT INTO usenmesgs (message, category) VALUES (?, ?)")
					if err != nil {
						fmt.Println(err)
					}
					statement.Exec(text, category)
				}
				db.Close()
			case "usen":
				db := dbConn()
				message := strings.Split(update.Message.Text, "/usen")
				message = strings.Split(message[1], " ")
				category := message[1]
				var ms []string
				rows, err := db.Query("SELECT message FROM usenmesgs WHERE category = ?", category)
				if err == nil {
					for rows.Next() {
						var m string
						if err := rows.Scan(&m); err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command :(")
							bot.Send(msg)
							break
						}
						ms = append(ms, m)
					}
				}
				
				//select random message from category
				if len(ms) > 0 {
					index := rand.Intn(len(ms))
					fmt.Println(ms[index])
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ms[index])
					bot.Send(msg)
					db.Close()
				}
			default: 
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command :(")
				bot.Send(msg)
			}
		}
	}
}