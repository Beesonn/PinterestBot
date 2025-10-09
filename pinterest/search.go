package pinterest

import (
	"fmt"
	"strings"

	"github.com/Mishel-07/PinterestBot/settings"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func FindImage(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.Message
	split := strings.SplitN(message.GetText(), " ", 2)
	if len(split) < 2 {
		message.Reply(b, "<b>No Query Provied So i can't send Photo, so Please Provide Query Eg:</b> <code>/pinterest Iron man</code>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return fmt.Errorf("no query provided")
	}

	query := split[1]
	msg, fck := message.Reply(b, "<b>Searching...🔎</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
	if fck != nil {
		fmt.Println(fck)
	}
	quotequery := strings.Replace(query, " ", "+", -1)
	urls, err := settings.SearchPinterest(quotequery)
	if err != nil {
		fmt.Println(err)
		message.Reply(b, "Image not found", &gotgbot.SendMessageOpts{})
		return err
	}

	media := make([]gotgbot.InputMedia, 0)
	count := 0
	for _, item := range urls.Data {
		if count == 10 {
			break
		}
		if item.URL != "" {
			media = append(media, gotgbot.InputMediaPhoto{Media: gotgbot.InputFileByURL(item.URL)})
		} else {
			fmt.Println("Skipped empty URL")
		}
		count++
	}
	if len(media) == 0 {
		message.Reply(b, "No Image found", &gotgbot.SendMessageOpts{})
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
		message.Reply(b, "No Image, please try again later.", &gotgbot.SendMessageOpts{})
		return err
	}

	return nil
}
