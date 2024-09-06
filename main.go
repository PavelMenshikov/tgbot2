package main

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var cardsPath = "./cards/"

func main() {

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не задан")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u) // Переместил переменную updates внутрь функции main

	rand.Seed(time.Now().UnixNano())

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Напиши /card, чтобы получить метафорическую карту.")
				bot.Send(msg)
			case "card":

				cardImage, err := getRandomCardImage(cardsPath)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось найти карты.")
					bot.Send(msg)
					continue
				}

				photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(cardImage))
				photo.Caption = "Твоя метафорическая карта"

				bot.Send(photo)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не знаю такой команды.")
				bot.Send(msg)
			}
		}
	}
}

func getRandomCardImage(path string) (string, error) {
	files, err := filepath.Glob(path + "*.*")
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", nil
	}

	return files[rand.Intn(len(files))], nil
}
