package main

import (
	"log"
	"strconv"

	"github.com/Syfaro/telegram-bot-api"
	"github.com/andygrunwald/go-jira"
)

func main() {
	jiraClient, err := jira.NewClient(nil, "https://your-instance.com")
	if err != nil {
		panic(err)
	}
	jiraClient.Authentication.SetBasicAuth("username", "password")

	//boards :=
	//issue, _, err := jiraClient.Issue.Get("SAM-33", nil)
	issue, _, err := jiraClient.Issue.Get("ISSUE-ID", nil)
	if err != nil {
		panic(err)
	}

	//boards := jiraClient.Board.GetBoard(boardID)
	//board := jiraClient.Issue.Get("SAM-33", nil)

	//fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)

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

			OriginalEstimate := strconv.Itoa(issue.Fields.TimeOriginalEstimate / 60)
			TimeSpent := strconv.Itoa(issue.Fields.TimeSpent / 60)
			TimeEstimate := strconv.Itoa(issue.Fields.TimeEstimate / 60)

			if Text == "/today" {
				log.Printf("[%s] %d %s", UserName, ChatID, Text)

				reply := issue.Fields.Summary + "(" + OriginalEstimate + "/" + TimeSpent + "/" + TimeEstimate + ")"

				msg := tgbotapi.NewMessage(ChatID, reply)

				bot.Send(msg)
			}
		}
	}
}
