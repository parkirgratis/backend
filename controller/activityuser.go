package controller

import (
    "context"
    "time"

    "github.com/gocroot/config"
    "github.com/gocroot/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func LogActivity(username, activity string) error {
    collection := config.Mongoconn.Collection("activity_logs")
    ctx := context.Background()

    log := model.LoginLog{
        ID:        primitive.NewObjectID().Hex(),
        Username:  username,
        Activity:  activity,
        Timestamp: time.Now(),
    }

    _, err := collection.InsertOne(ctx, log)
    if err != nil {
        return err
    }
    return nil
}
