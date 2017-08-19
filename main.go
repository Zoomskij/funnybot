package main
import (
  "github.com/Syfaro/telegram-bot-api"
  "log"
)

func main() {
  bot, err := tgbotapi.NewBotAPI("TOKEN")
  if err != nil {
    log.Panic(err)
  }

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

  //Initialize chanel
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upd, _ := bot.GetUpdatesChan(ucfg)

  //Reading the updates from chanel
	for {
		select {
		case update := <-upd:

			UserName := update.Message.From.UserName

			ChatID := update.Message.Chat.ID

			Text := update.Message.Text

			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			reply := "Hola Amigo!"

			msg := tgbotapi.NewMessage(ChatID, reply)

			bot.Send(msg)
		}
	}
}
