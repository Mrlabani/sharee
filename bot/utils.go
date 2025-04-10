package bot

import (
    "context"
    "fmt"
    "log"
    "os"
    "strings"
    "time"

    "filesharebot/models"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

func init() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    var err error
    client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
    if err != nil {
        log.Fatal(err)
    }
    collection = client.Database("filesharebot").Collection("files")
}

func IsAuthorizedUser(userID int64) bool {
    allowed := strings.Split(os.Getenv("AUTH_USERS"), ",")
    for _, id := range allowed {
        if fmt.Sprintf("%d", userID) == strings.TrimSpace(id) {
            return true
        }
    }
    return false
}

func HandleFileSave(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) tgbotapi.MessageConfig {
    file := models.File{
        FileID:    msg.Document.FileID,
        UniqueID:  msg.Document.FileUniqueID,
        FileName:  msg.Document.FileName,
        UserID:    msg.From.ID,
        Timestamp: time.Now(),
    }
    collection.InsertOne(context.Background(), file)
    baseURL := os.Getenv("WEB_BASE_URL")
    fileURL := fmt.Sprintf("%s/api/%s", baseURL, file.UniqueID)
    return tgbotapi.NewMessage(msg.Chat.ID, "Your file link: "+fileURL)
}

func HandleChannelEdit(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
    file := models.File{
        FileID:    msg.Document.FileID,
        UniqueID:  msg.Document.FileUniqueID,
        FileName:  msg.Document.FileName,
        UserID:    msg.Chat.ID,
        Timestamp: time.Now(),
    }
    collection.InsertOne(context.Background(), file)
    baseURL := os.Getenv("WEB_BASE_URL")
    fileURL := fmt.Sprintf("%s/api/%s", baseURL, file.UniqueID)
    edit := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, "File saved. [Access File]("+fileURL+")")
    edit.ParseMode = "Markdown"
    bot.Send(edit)
}