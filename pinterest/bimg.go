package pinterest

import (
	"fmt"
	"strings"

	"github.com/Mishel-07/PinterestBot/settings"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func BingImgCmd(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.Message
	split := strings.SplitN(message.GetText(), " ", 2)
	if len(split) < 2 {
		message.Reply(b, "<b>No Query Provied So i can't send Photo, so Please Provide Query</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return nil
	}

	query := split[1]
	msg, fck := message.Reply(b, "<b>Searching...🔎</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
	if fck != nil {
		return nil
	}
	urls, err := settings.SearchBing(query, 8)
	if err != nil {
		fmt.Println(err)
		message.Reply(b, "<b>Hey, No image Found So Report Here @XBOTSUPPORTS</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return err
	}

	media := make([]gotgbot.InputMedia, 0)
	for _, item := range urls {
		if item.URL != "" {
			media = append(media, gotgbot.InputMediaPhoto{
				Media: gotgbot.InputFileByURL(item.URL),
			})
		} else {
			fmt.Println("No image maybe Api not right now")
		}
	}

	if len(media) == 0 {
		message.Reply(b, "<b>Hey, No image Found So Report Here @XBOTSUPPORTS</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		b.DeleteMessage(msg.Chat.Id, msg.MessageId, &gotgbot.DeleteMessageOpts{})
		return fmt.Errorf("no valid media found to send")
	}

	_, err = b.SendMediaGroup(
		message.Chat.Id,
		media,
		&gotgbot.SendMediaGroupOpts{
			ReplyParameters: &gotgbot.ReplyParameters{
				MessageId: message.MessageId,
			},
		},
	)
	b.DeleteMessage(msg.Chat.Id, msg.MessageId, &gotgbot.DeleteMessageOpts{})
	if err != nil {
		fmt.Printf("Error sending media group: %s\n", err)
		message.Reply(b, "An error occurred, Report Here @XBOTSUPPORTS", &gotgbot.SendMessageOpts{})
		return err
	}

	return nil
}
