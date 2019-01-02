package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var userID string
				if event.Source != nil {
					userID = event.Source.UserID
				}

				switch {
				case strings.EqualFold(message.Text, "link"):
					//token link
					res, err := bot.IssueLinkToken(userID).Do()
					if err != nil {
						log.Println("Issue link token error, err=", err)
					}

					log.Println("Get user token:", res.LinkToken)

					if _, err = bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTemplateMessage("Account Link", linebot.NewButtonsTemplate(
							"",
							"account link",
							"account link",
							linebot.NewURIAction("Account Link", serverURL+"link?linkToken="+res.LinkToken)),
						),
					).Do(); err != nil {
						log.Println("err:", err)
						return
					}

					return
				case strings.EqualFold(message.Text, "list"):
					if _, err = bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTemplateMessage("List all user", linebot.NewButtonsTemplate(
							"",
							"List all cusotmers in provider website",
							"list all cusotmers",
							linebot.NewURIAction("List all cusotmers", serverURL)),
						),
					).Do(); err != nil {
						log.Println("err:", err)
						return
					}

					return
				}

				if _, err = bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("Currently your account is not linked to provider. \nThis is a starter for account link, check following actions.").
						WithQuickReplies(linebot.NewQuickReplyItems(
							linebot.NewQuickReplyButton(
								"",
								linebot.NewMessageAction("account link", "link")),
							linebot.NewQuickReplyButton(
								"",
								linebot.NewMessageAction("list user", "list")),
						)),
				).Do(); err != nil {
					log.Println("err:", err)
					return
				}
			}
		} else if event.Type == linebot.EventTypeAccountLink {
			// account link success.s
			log.Println("EventTypeAccountLink: source=", event.Source, " result=", event.AccountLink.Result)
		}
	}
}
