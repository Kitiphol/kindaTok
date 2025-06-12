package service

import (
    "github.com/RichardKnop/machinery/v2"
    "github.com/RichardKnop/machinery/v2/tasks"
)

func SendTaskToThumbnailCreator(server *machinery.Server, bucket, objectKey string) error {
    signature := &tasks.Signature{
        Name: "CreateThumbnail",
        Args: []tasks.Arg{
            {Type: "string", Value: bucket},
            {Type: "string", Value: objectKey},
        },
    }
    _, err := server.SendTask(signature)
    return err
}
