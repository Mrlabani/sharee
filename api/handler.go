package api

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"

    "filesharebot/models"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
    collection = client.Database("filesharebot").Collection("files")
}

func Handler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/api/"):]

    var file models.File
    err := collection.FindOne(context.Background(), map[string]string{"unique_id": id}).Decode(&file)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    res := map[string]string{
        "file_id":   file.FileID,
        "file_name": file.FileName,
        "telegram":  fmt.Sprintf("https://t.me/%s?start=%s", os.Getenv("BOT_USERNAME"), file.UniqueID),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}
