package pinterest

import (
	"fmt"
	"github.com/Mishel-07/PinterestBot/settings"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"regexp"
	"strings"
)

func ExtractURL(message string) string {
	pattern := regexp.MustCompile(`https?://\S+`)
	match := pattern.FindString(message)
	return match
}

func DownloadSend(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	if message.Chat.Type != "private" {
		return nil
	}
	chk := message.Text
	if strings.HasPrefix(chk, "/") {
		return nil
	}
	pattern := regexp.MustCompile(`https://pin\.it/(?P<url>[\w]+)`)
	if !pattern.MatchString(chk) {
		return nil
	}

	link := ExtractURL(chk)
	url, err := settings.PinterestDownload(link)
	fmt.Println(url)
	if err != nil {
		message.Reply(b, "opps! An Error Occured Report on @XBOTSUPPORTS", &gotgbot.SendMessageOpts{})
		fmt.Println(err)
		return err
	}
	photo := gotgbot.InputMediaPhoto{
		Media: gotgbot.InputFileByURL(url),
	}
	_, uploadErr := b.SendPhoto(ctx.EffectiveChat.Id, photo.Media, &gotgbot.SendPhotoOpts{})
	if uploadErr != nil {
		message.Reply(b, "Failed to Send Photo", &gotgbot.SendMessageOpts{})
		return err
	}
	return nil
}
