package chat

import (
	"fmt"
	"os"
	"testing"
)

func TestModel(t *testing.T) {
	Debug = true
	options := DefaultGiteeAIOptions
	options.APIKey = os.Getenv(ENV_AI_API_KEY)
	ai := NewGiteeAI(options, false)
	fmt.Println(string(ai.MetaData().ToJSON()))
	ai.SetOnReceive(func(message *TResponse) {
		if message != nil {
			choices := message.Choices
			for _, choice := range choices {
				print(choice.Delta.Content)
			}
		}
	})
	//ai.Chat("你好")
	ai.ChatStream("你好")
}
