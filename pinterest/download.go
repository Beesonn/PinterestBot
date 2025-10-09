package pinterest

import (
	"regexp"
	"strings"

	"github.com/Mishel-07/PinterestBot/settings"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func ExtractURL(message string) string {
	pattern := regexp.MustCompile(`https?://\S+`)
	match := pattern.FindString(message)
	return match
}

func DownloadSend(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage	
	chk := message.Text
	if strings.HasPrefix(chk, "/") {
		return nil
	}
	pattern := regexp.MustCompile(`https://pin\.it/(?P<url>[\w]+)`)
	if !pattern.MatchString(chk) {
		return nil
	}

	link := ExtractURL(chk)
	url, err := settings.DownloadPinterestImage(link)
	if err != nil {
		message.Reply(b, "opps! An Error Occured Report on @XBOTSUPPORTS", nil)		
		return err
	}
	photo := gotgbot.InputMediaPhoto{
		Media: gotgbot.InputFileByURL(url),
	}
	_, uploadErr := b.SendPhoto(ctx.EffectiveChat.Id, photo.Media, &gotgbot.SendPhotoOpts{ReplyParameters: &gotgbot.ReplyParameters{MessageId: message.MessageId},})
	if uploadErr != nil {
		message.Reply(b, "Failed to Send Photo", nil)
		return err
	}
	return nil
}
