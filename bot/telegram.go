package bot

import (
    "filesharebot/utils"
    "log"
    "os"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot() {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
    if err != nil {
        log.Panic(err)
    }

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        log.Println("Received update...")

        if update.Message != nil {
            log.Printf("From: %s | Text: %s", update.Message.From.UserName, update.Message.Text)
            userID := update.Message.From.ID

            if !utils.IsAuthorizedUser(userID) {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Not authorized."))
                continue
            }

            if update.Message.Chat.IsPrivate() {
                if update.Message.Document != nil || update.Message.Video != nil ||
                    update.Message.Audio != nil || update.Message.Photo != nil {
                    msg := utils.HandleFileSave(bot, update.Message)
                    bot.Send(msg)
                } else {
                    bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please send a file to get a sharable link."))
                }
            }
        }

        if update.ChannelPost != nil {
            log.Println("Channel post received")
            if update.ChannelPost.Document != nil || update.ChannelPost.Video != nil || update.ChannelPost.Photo != nil {
                utils.HandleChannelEdit(bot, update.ChannelPost)
            }
        }
    }
}