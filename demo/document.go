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
			println()
		}
	})
	ai.SetOnFail(func(message *chat.TResponseError) {
		fmt.Println("fail:", message.Code, message.Message, message.Type)
	})
	wd, _ := os.Getwd()
	sourcePath := filepath.Join(wd, "chat")
	demoPath := filepath.Join(wd, "demo")

	sourceCode := tool.LoadSourceFile(sourcePath, false)
	println("源码:", len(sourceCode))
	demoCode := tool.LoadSourceFile(demoPath, false)
	println("示例:", len(demoCode))

	contentBuf := bytes.Buffer{}
	contentBuf.WriteString("你现在为这个项目编写使用说明文档，根据用户要求编写项目介绍的 README.md 文档。\n")
	contentBuf.WriteString("\n源代码:\n")
	contentBuf.Write(sourceCode)
	contentBuf.WriteString("\n参考示例:\n")
	contentBuf.Write(demoCode)
	contentBuf.WriteString("\n")
	println("总:", contentBuf.Len())

	time.Sleep(time.Second / 2)

	ai.ChatStream(contentBuf.String())
	println()
}
