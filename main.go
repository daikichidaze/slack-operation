package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// アクセストークンを使用してクライアントを生成する

	flag.Parse()
	tkn := flag.Arg(0)
	api := slack.New(tkn)
	_, chn_id := get_calendar_channel(api)

	prm := slack.GetConversationHistoryParameters{ChannelID: chn_id}
	res, err := api.GetConversationHistory(&prm)
	check(err)

	delete_title := "た：みーてぃんぐ"

	for _, m := range res.Messages {
		for _, att := range m.Attachments {
			if att.Title == delete_title {
				a, b, err := api.DeleteMessage(chn_id, m.Timestamp)
				check(err)
				fmt.Println(a, b)
			}
		}
	}

}

func get_calendar_channel(api *slack.Client) (*slack.Channel, string) {
	var result string
	calender_name := "calendar"
	prm := slack.GetConversationsParameters{}

	channels, _, err := api.GetConversations(&prm)
	check(err)

	for _, c := range channels {

		if strings.Contains(c.Name, calender_name) {
			result = c.ID
		}
	}

	c, err := api.GetConversationInfo(result, true)
	check(err)

	return c, result

}
