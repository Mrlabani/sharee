package models

import "time"

type File struct {
    FileID    string    `bson:"file_id"`
    UniqueID  string    `bson:"unique_id"`
    FileName  string    `bson:"file_name"`
    UserID    int64     `bson:"user_id"`
    Timestamp time.Time `bson:"timestamp"`
}
