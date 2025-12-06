package settings

import (
    "fmt"

    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func Start(b *gotgbot.Bot, ctx *ext.Context) error {
    message := ctx.Message
    if message.Chat.Type != "private" {
        return nil
    }
    buttons := gotgbot.InlineKeyboardMarkup{
        InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
            {
                {Text: "Source", Url: "https://github.com/Beesonn/PinterestBot"},
            },
            {
                {Text: "Support", Url: "https://t.me/XBOTSUPPORTS"},
            },
        },
    }
    _, err := message.Reply(b, "<b>Hey, I'm PinterestBot. You can search for Pinterest videos or photos, and you can download them too. I can also provide Google images, Bing images, etc.</b>\n\n<b>Commands:</b>\n\n/pinterest - search and download pinterest image\n/wallpaper - wallpaper search\n/img - from bing image\n\n<b>Tools:</b>\n\nSend me a Pinterest url I'll give that photo/video", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML, ReplyMarkup: buttons})
    if err != nil {
        fmt.Println(err)
    }
    return nil
}
