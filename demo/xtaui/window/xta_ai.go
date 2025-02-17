package window

import (
	"bytes"
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/xta/chat"
	"strings"
)

func (m *TMainWindow) initXTASDK() {
	options := chat.DefaultGiteeAIOptions
	options.APIKey = "FONF7P9SWBWG0DDSYVPJJKHH3WTAAW2VMBS8YV1O" // os.Getenv(chat.ENV_AI_API_KEY)
	m.ai = chat.NewGiteeAI(options, false)
	m.ai.System("【系统角色】你具备跨领域知识整合与结构化推理能力的智能助手。始终遵循：事实准确性 > 响应速度 > 表达流畅度的优先级原则。【响应规范】1. 解析阶段：识别问题类型（事实/观点/方法需求），标注关键信息置信度2. 处理阶段：- 事实类：提供最新权威信源+时间戳- 观点类：多视角分析+概率评估 - 方法类：分步实施框架+风险预案3. 输出阶段：采用「结论-依据-延伸」结构，技术概念附带白话解释【安全协议】- 对潜在争议内容自动附加免责声明- 医疗/法律建议必须提示咨询专业人士- 实时监测对话情感倾向，对焦虑/紧急表达启动安抚话术")
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
		m.message.Lines().Add("模型: " + m.ai.Model())
		//m.message.Lines().Add("APIKEY: ........." + m.ai.APIKey()[5:10] + "............")
		m.SetCaption(m.title + " " + m.ai.Model())
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
			if fw.isSend {
				continue
			}
			fw.isSend = true
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
