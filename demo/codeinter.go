package main

import (
	"bytes"
	"fmt"
	"github.com/energye/xta/chat"
	"github.com/energye/xta/demo/tool"
	"os"
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
	// 自定义一些资源
	sourceCode := tool.LoadSourceFile("E:\\SWT\\gopath\\src\\github.com\\energye\\energy\\cef", true)
	println("源码:", len(sourceCode))
	doc := tool.LoadFile("E:\\SWT\\gopath\\src\\github.com\\energye\\energye.github.io\\zh\\course")
	println("文档:", len(doc))
	demoCode := tool.LoadSourceFile("E:\\SWT\\gopath\\src\\github.com\\energye\\energy\\examples\\many-browser", false)
	println("示例:", len(demoCode))

	contentBuf := bytes.Buffer{}
	contentBuf.WriteString("你主要负责为这个示例解读，根据用户要求解读示例代码并生成章节内容。\n")
	//contentBuf.WriteString("\n参考内容:\n")
	//contentBuf.Write(doc)
	contentBuf.WriteString("\n代码内容:\n")
	contentBuf.Write(sourceCode)
	contentBuf.WriteString("\n代码示例:\n")
	contentBuf.Write(demoCode)
	contentBuf.WriteString("\n")
	println("总:", contentBuf.Len())

	time.Sleep(time.Second / 2)

	ai.ChatStream(contentBuf.String())
	println()
}
