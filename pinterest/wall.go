package pinterest

import (
    "fmt"
    "strings"

    "github.com/Mishel-07/PinterestBot/settings"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func WallSearch(b *gotgbot.Bot, ctx *ext.Context) error {
    message := ctx.Message
    split := strings.SplitN(message.GetText(), " ", 2)
    if len(split) < 2 {
        message.Reply(b, "<b>No Query Provided So I can't send Photo, so Please Provide Query</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
        return fmt.Errorf("no query provided")
    }

    query := split[1]
    msg, fck := message.Reply(b, "<b>Searching...🔎</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
    if fck != nil {
        return nil
    }
    quotequery := strings.Replace(query, " ", "+", -1)
    images := settings.ScrapWallpapers(quotequery)

    media := make([]gotgbot.InputMedia, 0)
    count := 0
    for _, item := range images {
        if count == 10 {
            break
        }
        media = append(media, gotgbot.InputMediaPhoto{
            Media: gotgbot.InputFileByURL(item),
        })
        count++
    }

    if len(media) == 0 {
        message.Reply(b, "No Image found", &gotgbot.SendMessageOpts{})
        b.DeleteMessage(msg.Chat.Id, msg.MessageId, &gotgbot.DeleteMessageOpts{})
        return fmt.Errorf("no valid media found to send")
    }

    b.SendMediaGroup(
        message.Chat.Id,
        media,
        &gotgbot.SendMediaGroupOpts{
            ReplyParameters: &gotgbot.ReplyParameters{
                MessageId: message.MessageId,
            },
        },
    )
    b.DeleteMessage(msg.Chat.Id, msg.MessageId, &gotgbot.DeleteMessageOpts{})

    return nil
}
