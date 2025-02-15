package window

import (
	"bytes"
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/xta/chat"
	"os"
	"strings"
)

func (m *TMainWindow) initXTASDK() {
	options := chat.DefaultGiteeAIOptions
	options.APIKey = os.Getenv(chat.ENV_AI_API_KEY)
	m.ai = chat.NewGiteeAI(options, false)

	isFirstRec := false
	m.ai.SetOnReceive(func(message *chat.TResponse) {
		if !isFirstRec {
			m.message.Lines().Add("回复: " + nowDatetime())
			isFirstRec = true
		}
		// 在异步UI线程里操作
		lcl.RunOnMainThreadAsync(func(id uint32) {
			if message != nil {
				if message.Error != "" {
					s := fmt.Sprintf("错误: %v %v", message.Error, message.ErrorType)
					m.message.Lines().Add(s)
					if m.saveFileBuf != nil {
						m.saveFileBuf.WriteString(s)
						m.saveFileBuf.Flush()
					}
				}
				choices := message.Choices
				for _, choice := range choices {
					if strings.Contains(choice.Delta.Content, "\n") {
						m.message.Lines().Add(choice.Delta.Content)
					} else {
						m.message.SetSelStart(int32(len(m.message.Lines().Text())))
						m.message.SetSelText(choice.Delta.Content)
					}
					if m.saveFileBuf != nil {
						m.saveFileBuf.WriteString(choice.Delta.Content)
						m.saveFileBuf.Flush()
					}
				}
			} else {
				fmt.Println("结束")
				m.message.Lines().Add("")
				m.chatBtn.SetEnabled(true)
				isFirstRec = false
			}
		})
	})
	m.ai.SetOnFail(func(message *chat.TResponseError) {
		lcl.RunOnMainThreadSync(func() {
			s := fmt.Sprintf("  错误: %v %v %v", message.Code, message.Message, message.Type)
			m.message.Lines().Add(s)
			m.chatBtn.SetEnabled(true)
			isFirstRec = false
		})
	})
	lcl.RunOnMainThreadSync(func() {
		m.message.Lines().Add("XTA - AI SDK 初始化完成")
		m.message.Lines().Add("模型: " + m.ai.Name())
		//m.message.Lines().Add("APIKEY: ........." + m.ai.APIKey()[5:10] + "............")
		m.SetCaption(m.title + " " + m.ai.Name())
	})
}

func (m *TMainWindow) SendMessage(sender lcl.IObject) {
	msg := m.chat.Lines().Text()
	if msg != "" {
		m.message.Lines().Add("我: " + nowDatetime())
		m.message.Lines().Add("  " + msg)
		buf := bytes.Buffer{}
		buf.WriteString(msg + "\n")
		for _, fw := range m.fileWindow {
			buf.WriteString(fw.text.Text() + "\n")
			buf.WriteString(strings.Join(fw.fileContent, "\n"))
		}
		m.sendMessage(buf.String())
		m.chat.SetText("")
	}
}

func (m *TMainWindow) sendMessage(content string) {
	// 在协程里操作
	go m.ai.ChatStream(content)
	m.chatBtn.SetEnabled(false)
}
