package main

import (
	"bytes"
	"fmt"
	"github.com/energye/xta/chat"
	"github.com/energye/xta/demo/tool"
	"os"
	"path/filepath"
	"time"
)

func main() {
	chat.Debug = true

	options := chat.DefaultGiteeAIOptions
	options.APIKey = os.Getenv(chat.ENV_AI_API_KEY)
	ai := chat.NewGiteeAI(options, false)

	ai.Header().Set("x-failover-enabled", "true")
	ai.Header().Set("x-package", "1910")
	ai.Header().Set("x-pipeline-tag", "conversational")
	ai.Header().Set("x-trial-enabled", "true")

	ai.SetOnReceive(func(message *chat.TResponse) {
		if message != nil {
			if message.Error != "" {
				fmt.Println(message.Error, message.ErrorType)
			} else {
				choices := message.Choices
				for _, choice := range choices {
					fmt.Print(choice.Delta.Content)
				}
			}
		} else {
			println("结束")
		}
	})
	ai.SetOnFail(func(message *chat.TResponseError) {
		fmt.Println("fail:", message.Code, message.Message, message.Type)
	})

	ai.System("你现在是一个文档编写专家, 根据用户要求编写文档.")

	wd, _ := os.Getwd()
	docPath := filepath.Join(wd, "demo", "data")

	docCode := tool.LoadFile(docPath)
	println("文档:", len(docCode))

	question := bytes.Buffer{}
	question.WriteString("你现在为这个项目编写使用说明文档，根据用户要求编写项目介绍的文档。\n")
	question.WriteString("\n参考文档:\n")
	question.Write(docCode)
	question.WriteString("\n")
	println("总:", question.Len())

	time.Sleep(time.Second / 2)

	ai.ChatStream(question.String())
	println()
}
