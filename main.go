package main

import (
	"log"
	"strconv"
	"strings"

	"./hubstaff-wrapper"

	"github.com/Syfaro/telegram-bot-api"
	"github.com/andygrunwald/go-jira"
)

const START_TASKS = 59
const COUNT_TASKS = 66

func main() {
	jiraClient, err := jira.NewClient(nil, "https://INSTANCE.atlassian.net")
	if err != nil {
		panic(err)
	}
	jiraClient.Authentication.SetBasicAuth("EMAIL", "PASS")

	bot, err := tgbotapi.NewBotAPI("TOkEN")
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
					"/auth %EMAIL% %PASSWORD% \n" +
					"/logout \n" +
					"/today - return one issue (static issue)\n" +
					"/getall - return all issues \n" +
					"/upper - show with excess estimate \n" +
					"/less - show with lack of estimate "
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			}

			if Text == "/today" {
				if jiraClient.Authentication.Authenticated() == false {
					msg := tgbotapi.NewMessage(ChatID, "You're not authorized!")
					bot.Send(msg)
					break
				}

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
				for i := START_TASKS; i <= COUNT_TASKS; i++ {
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
				for i := START_TASKS; i <= COUNT_TASKS; i++ {
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
				for i := START_TASKS; i <= COUNT_TASKS; i++ {
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
				//reply := hubwrapper.Auth() + "\n"
				//reply += hubwrapper.GetUsers()
				reply := hubwrapper.GetTasks()
				//reply := hubwrapper.auth()
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			}

			if Text == "/boards" {
				var reply string
				var boards *jira.BoardsList
				boards, _, err = jiraClient.Board.GetAllBoards(nil)
				if err != nil {
					panic(err)
				}

				for _, board := range boards.Values {
					reply += board.Name + "\n"
				}
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			}

			// if Text == "/projects" {
			// 	var reply string
			// 	var projects *jira.ProjectList
			// 	projects, _, err = jiraClient.Project.GetList()
			// 	if err != nil {
			// 		panic(err)
			// 	}
			//
			// 	for _, project := range projects {
			// 		reply += project.Key + "\n"
			// 	}
			// 	msg := tgbotapi.NewMessage(ChatID, reply)
			// 	bot.Send(msg)
			// }

			// if Text == "/current-user" {
			// 	var reply string
			// 	session := jiraClient.Authentication.
			// 	currentUser, _, err := jiraClient.Authentication.GetCurrentUser()
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// 	reply = currentUser.Name
			// 	msg := tgbotapi.NewMessage(ChatID, reply)
			// 	bot.Send(msg)
			// }

			if strings.Contains(Text, "/auth") {
				split := strings.Split(Text, " ")
				email := split[1]
				password := split[2]
				if jiraClient.Authentication.Authenticated() == true {
					err = jiraClient.Authentication.Logout()
					if err != nil {
						panic(err)
					}
				}

				jiraClient.Authentication.AcquireSessionCookie(email, password)

				if jiraClient.Authentication.Authenticated() == true {
					msg := tgbotapi.NewMessage(ChatID, "Auth Successfuly")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(ChatID, "Auth Unfortunately")
					bot.Send(msg)
				}
			}

			if Text == "/logout" {
				if jiraClient.Authentication.Authenticated() == true {
					err = jiraClient.Authentication.Logout()
					msg := tgbotapi.NewMessage(ChatID, "Logout Successfuly")
					bot.Send(msg)
					if err != nil {
						panic(err)
					}
				} else {
					msg := tgbotapi.NewMessage(ChatID, "Auth Unfortunately")
					bot.Send(msg)
				}
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
		reply += " (" + issue.Fields.Assignee.DisplayName + ")"
	}
	reply += "\n"

	var indexSpent int
	if issue.Fields.TimeOriginalEstimate > issue.Fields.TimeSpent {
		indexSpent = issue.Fields.TimeSpent / 3600
	} else {
		indexSpent = issue.Fields.TimeSpent/3600 - (issue.Fields.TimeSpent/3600 - issue.Fields.TimeOriginalEstimate/3600)
	}
	for i := 0; i < indexSpent; i++ {
		reply += "\xF0\x9F\x94\xB5"
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
