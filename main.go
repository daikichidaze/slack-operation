package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// arg 0: slack access token
// arg 1: channel name
// arg 2: delete message title
func main() {

	// Get access token from args
	flag.Parse()
	tkn := flag.Arg(0)
	api := slack.New(tkn)

	channel_name := flag.Arg(1)                       // Target name of channel
	chn_id := get_calendar_channel(api, channel_name) // Get target channel id

	// Get message list
	prm := slack.GetConversationHistoryParameters{ChannelID: chn_id, Limit: 1000}
	res, err := api.GetConversationHistory(&prm)
	check(err)

	delete_title := flag.Arg(2) // Target title name to delete
	var cnt int

	for _, m := range res.Messages {
		// For google calendar
		for _, att := range m.Attachments {
			if att.Title == delete_title {
				cnt = delete_message(api, chn_id, m.Timestamp, cnt)

				if cnt%25 == 0 {
					fmt.Println(time.Now(), cnt, "messages were deleted")
				}
			}
		}

		//For tobuy
		// if strings.Contains(m.Text, delete_title) {
		// 	cnt = delete_message(api, chn_id, m.Timestamp, cnt)
		// }
	}
	fmt.Println(cnt, "messages were deleted")

}

func get_calendar_channel(api *slack.Client, channel_name string) string {
	var result string

	prm := slack.GetConversationsParameters{}
	channels, _, err := api.GetConversations(&prm)
	check(err)

	// Find channel
	for _, c := range channels {
		if strings.Contains(c.Name, channel_name) {
			result = c.ID
		}
	}

	return result
}

func delete_message(api *slack.Client, chn_id string, timestamp string, cnt int) int {
	_, _, err := api.DeleteMessage(chn_id, timestamp)
	check(err)
	time.Sleep(750 * time.Millisecond) // To escape the api limit
	cnt++
	return cnt
}
