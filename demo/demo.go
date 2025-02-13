package main

import (
	"fmt"
	"github.com/energye/xta/chat"
	"os"
)

func main() {
	chat.Debug = false
	options := chat.DefaultGiteeAIOptions
	options.APIKey = os.Getenv(chat.ENV_AI_API_KEY)
	ai := chat.NewGiteeAI(options, false)
	chatBox := func() {
		fmt.Println(">")
		var input string
		fmt.Scan(&input)
		if input != "" {
			ai.ChatStream(input)
		}
	}
	ai.SetOnReceive(func(message *chat.TResponse) {
		if message != nil {
			if message.Error != "" {
				fmt.Println(message.Error, message.ErrorType)
			}
			choices := message.Choices
			for _, choice := range choices {
				fmt.Print(choice.Delta.Content)
			}
		} else {
			println()
			chatBox()
		}
	})
	ai.SetOnFail(func(message *chat.TResponseError) {
		fmt.Println("fail:", message.Code, message.Message, message.Type)
	})
	chatBox()
}
