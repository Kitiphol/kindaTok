package machineryutil

import (
    "github.com/RichardKnop/machinery/v2/tasks"
)

func SendChunkVideoTask(input string) error {
    
    server, err := CreateMachineryServer()
    if err != nil {
        return err
    }
    signature := &tasks.Signature{
        Name: "ChunkVideo",
        Args: []tasks.Arg{
            {
                Type:  "string",
                Value: input,
            },
        },
    }
    _, err = server.SendTask(signature)
    return err
}