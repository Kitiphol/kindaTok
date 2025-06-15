package machineryUtil

import (
    "github.com/RichardKnop/machinery/v2"
    "github.com/RichardKnop/machinery/v2/tasks"
)

func SendTaskToConvertor(server *machinery.Server, taskName string, args ...tasks.Arg) error {
    signature := &tasks.Signature{
        Name: taskName,
        Args: args,
    }
    _, err := server.SendTask(signature)
    return err
}