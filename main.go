package main

import (
	"log"
	"strconv"

	"github.com/Syfaro/telegram-bot-api"
	"github.com/andygrunwald/go-jira"
)

func main() {
	jiraClient, err := jira.NewClient(nil, "https://INSTANCE.COM")
	if err != nil {
		panic(err)
	}
	jiraClient.Authentication.SetBasicAuth("EMAIL", "PASSWORD")

	bot, err := tgbotapi.NewBotAPI("TELEGRAM.TOKEN")
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

			ChatID := update.Message.Chat.ID

			Text := update.Message.Text

			if Text == "/help" {
				reply := "Current commands: \n" +
					"/help - current help \n" +
					"/today - return one issue (static issue)\n" +
					"/getall - return all issues \n" +
					"/upper - show with excess estimate \n" +
					"/less - show with lack of estimate "
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			}

			if Text == "/today" {
				var reply string
				issue, _, err := jiraClient.Issue.Get("SAM-9", nil)
				if err != nil {
					panic(err)
				}

				reply += getReply(issue)

				msg := tgbotapi.NewMessage(ChatID, reply)

				bot.Send(msg)
			}

			if Text == "/getall" {
				var reply string
				for i := 1; i <= 39; i++ {
					issueId := "SAM-" + strconv.Itoa(i)
					issue, _, err := jiraClient.Issue.Get(issueId, nil)
					if err != nil {
						panic(err)
					}

					reply += getReply(issue)
				}
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			}

			if Text == "/upper" {
				var reply string
				for i := 1; i <= 39; i++ {
					issueId := "SAM-" + strconv.Itoa(i)
					issue, _, err := jiraClient.Issue.Get(issueId, nil)
					if err != nil {
						panic(err)
					}
					if (issue.Fields.TimeOriginalEstimate + 10*60) < issue.Fields.TimeSpent {
						reply += getReply(issue)
					}
				}
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			}

			if Text == "/less" {
				var reply string
				for i := 1; i <= 39; i++ {
					issueId := "SAM-" + strconv.Itoa(i)
					issue, _, err := jiraClient.Issue.Get(issueId, nil)
					if err != nil {
						panic(err)
					}
					if (issue.Fields.TimeOriginalEstimate) > (issue.Fields.TimeSpent + 10*60) {
						reply += getReply(issue)
					}
				}
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			}

			if Text == "/hub" {
				reply := hubauth()
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			}
		}
	}

}

func getReply(issue *jira.Issue) string {
	var reply string
	OriginalEstimate := strconv.Itoa(issue.Fields.TimeOriginalEstimate/3600) + "h" + strconv.Itoa(issue.Fields.TimeOriginalEstimate/60%60) + "m"
	TimeSpent := strconv.Itoa(issue.Fields.TimeSpent/3600) + "h" + strconv.Itoa(issue.Fields.TimeSpent/60%60) + "m"
	TimeEstimate := strconv.Itoa(issue.Fields.TimeEstimate/3600) + "h" + strconv.Itoa(issue.Fields.TimeEstimate/60%60) + "m"
	reply += issue.Key + " " + issue.Fields.Summary + "\n"
	reply += " ([" + OriginalEstimate + "]  [" + TimeSpent + "]  [" + TimeEstimate + "])"
	if issue.Fields.Assignee != nil {
		reply += " (" + issue.Fields.Assignee.DisplayName + ")\n"
	}

	var indexSpent int
	if issue.Fields.TimeOriginalEstimate > issue.Fields.TimeSpent {
		indexSpent = issue.Fields.TimeSpent / 3600
	} else {
		indexSpent = issue.Fields.TimeSpent/3600 - (issue.Fields.TimeSpent/3600 - issue.Fields.TimeOriginalEstimate/3600)
	}
	for i := 0; i < indexSpent; i++ {
		reply += "\xE2\x9A\xAB"
	}
	for i := 0; i < issue.Fields.TimeOriginalEstimate/3600-issue.Fields.TimeSpent/3600; i++ {
		reply += "\xE2\x9A\xAA"
	}
	for i := 0; i < issue.Fields.TimeSpent/3600-issue.Fields.TimeOriginalEstimate/3600; i++ {
		reply += "\xF0\x9F\x94\xB4"
	}
	reply += "\n\n"
	return reply
}
