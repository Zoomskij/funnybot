package main
import (
  "github.com/Syfaro/telegram-bot-api"
  "github.com/andygrunwald/go-jira"
  "log"
  "fmt"
)

func main() {
  jiraClient, err := jira.NewClient(nil, "https://jira.instance.com")
	if err != nil {
		panic(err)
	}
	jiraClient.Authentication.SetBasicAuth("username", "passowrd")

	issue, _, err := jiraClient.Issue.Get("TASK-KEY", nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)

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

			reply := issue.Fields.Summary

			msg := tgbotapi.NewMessage(ChatID, reply)

			bot.Send(msg)
		}
	}
}
