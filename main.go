package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/Mishel-07/PinterestBot/pinterest"
    "github.com/Mishel-07/PinterestBot/settings"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/inlinequery"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func KeepOnline(url string) {
    for {
        resp, err := http.Get(url)
        if err != nil {
            fmt.Println("Error:", err)
        }
        defer resp.Body.Close()
        time.Sleep(41 * time.Second)
    }
}

func main() {
    token := os.Getenv("TOKEN")
    if token == "" {
        panic("TOKEN environment variable is empty")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    webhook := os.Getenv("WEBHOOK")
    if webhook == "" {
        webhook = "true"
    }
    b, err := gotgbot.NewBot(token, nil)
    if err != nil {
        panic("failed to create new bot: " + err.Error())
    }
    if webhook != "false" {
        go func() {
            http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
                fmt.Fprintf(w, "Hello World")
            })

            http.ListenAndServe(":"+port, nil)
        }()
        url := os.Getenv("URL")
        if url != "" {
            go KeepOnline(url)
        }
    }

    dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
        Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
            log.Println("an error occurred while handling update:", err.Error())
            return ext.DispatcherActionNoop
        },
        MaxRoutines: ext.DefaultMaxRoutines,
    })
    updater := ext.NewUpdater(dispatcher, nil)
    dispatcher.AddHandler(handlers.NewCommand("start", settings.Start))
    dispatcher.AddHandler(handlers.NewCommand("pinterest", pinterest.FindImage))
    dispatcher.AddHandler(handlers.NewCommand("wallpaper", pinterest.WallSearch))
    dispatcher.AddHandler(handlers.NewCommand("img", pinterest.BingImgCmd))
    dispatcher.AddHandler(handlers.NewMessage(message.Text, pinterest.DownloadSend))
    dispatcher.AddHandler(handlers.NewInlineQuery(inlinequery.All, pinterest.FindImageInline))

    err = updater.StartPolling(b, &ext.PollingOpts{
        DropPendingUpdates: true,
        GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
            Timeout: 9,
            RequestOpts: &gotgbot.RequestOpts{
                Timeout: time.Second * 10,
            },
        },
    })
    if err != nil {
        panic("failed to start polling: " + err.Error())
    }
    log.Printf("%s has been started...\n", b.User.Username)

    updater.Idle()
}
