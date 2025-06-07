package machineryutil

import (
    "github.com/RichardKnop/machinery/v2/tasks"
)


func SendExampleMessage() error {
    server, err := CreateMachineryServer()
    if err != nil {
        return err
    }

    // Define the task signature (replace "YourTaskName" and args as needed)
    signature := &tasks.Signature{
        Name: "YourTaskName",
        Args: []tasks.Arg{
            {
                Type:  "string",
                Value: "Hello, World!",
            },
        },
    }

    _, err = server.SendTask(signature)
    return err
}