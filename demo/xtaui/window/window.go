package window

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/xta/chat"
	"os"
	"strings"
	"time"
	_ "xta/xtaui/syso"
)

type TMainWindow struct {
	lcl.TForm
	message lcl.IMemo
	chat    lcl.IMemo

	ai      chat.IGiteeAI
	chatBtn lcl.IButton
}

var MainWindow TMainWindow

func (m *TMainWindow) FormCreate(sender lcl.IObject) {
	m.SetCaption("ENERGY - XTA Chat UI")
	m.SetPosition(types.PoScreenCenter)
	m.SetWidth(800)
	m.SetHeight(600)

	png := lcl.NewPngImage()
	png.LoadFromFSFile("assets/icon.png")
	lcl.Application.Icon().Assign(png)
	png.Free()

	m.Init()
}

func (m *TMainWindow) Init() {
	go m.initXTA()

	// 消息
	m.message = lcl.NewMemo(m)
	m.message.SetParent(m)
	m.message.SetLeft(150)
	m.message.SetWidth(m.Width() - 150)
	m.message.SetHeight(m.Height() - 150)
	m.message.SetBorderStyle(types.BsNone)
	m.message.SetReadOnly(true)
	m.message.SetScrollBars(types.SsAutoBoth)
	m.message.SetWordWrap(true)
	m.message.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight, types.AkBottom))

	// 聊天
	chatLabel := lcl.NewLabel(m)
	chatLabel.SetParent(m)
	chatLabel.SetCaption("发送消息")
	chatLabel.SetTop(m.message.Height() + 40)
	chatLabel.SetLeft(20)
	chatLabel.Font().SetSize(18)
	chatLabel.Font().SetColor(colors.ClGray)
	chatLabel.SetAnchors(types.NewSet(types.AkLeft, types.AkBottom))

	m.chat = lcl.NewMemo(m)
	m.chat.SetParent(m)
	m.chat.SetBorderStyle(types.BsNone)
	m.chat.SetScrollBars(types.SsAutoBoth)
	m.chat.SetTop(m.message.Height() + 3)
	m.chat.SetLeft(150)
	m.chat.SetWidth(m.Width())
	m.chat.SetHeight(100)
	m.chat.SetAnchors(types.NewSet(types.AkLeft, types.AkRight, types.AkBottom))

	// 发送消息
	m.chatBtn = lcl.NewButton(m)
	m.chatBtn.SetParent(m)
	m.chatBtn.SetTop(m.chat.Top() + m.chat.Height() + 3)
	m.chatBtn.SetWidth(100)
	m.chatBtn.SetHeight(45)
	m.chatBtn.SetCaption("发 送")
	m.chatBtn.Font().SetSize(16)
	m.chatBtn.SetLeft(m.Width() - 100)
	m.chatBtn.SetAnchors(types.NewSet(types.AkRight, types.AkBottom))
	m.chatBtn.SetOnClick(m.SendMessage)

	m.SetOnShow(func(sender lcl.IObject) {
		m.chat.SetFocus()
	})
}

func nowDatetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (m *TMainWindow) initXTA() {
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
				}
				choices := message.Choices
				for _, choice := range choices {
					if strings.Contains(choice.Delta.Content, "\n") {
						m.message.Lines().Add(choice.Delta.Content)
					} else {
						m.message.SetSelStart(int32(len(m.message.Lines().Text())))
						m.message.SetSelText(choice.Delta.Content)
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
}

func (m *TMainWindow) SendMessage(sender lcl.IObject) {
	msg := m.chat.Lines().Text()
	if msg != "" {
		m.message.Lines().Add("我: " + nowDatetime())
		m.message.Lines().Add("  " + msg)
		// 在协程里操作
		go m.ai.ChatStream(msg)
		m.chat.SetText("")
		m.chatBtn.SetEnabled(false)
	}
}
