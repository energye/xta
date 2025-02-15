package window

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/xta/chat"
	"io/fs"
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

	selFileBtn lcl.IButton
	selDirDlg  lcl.ISelectDirectoryDialog

	saveChatBtn lcl.IButton
	saveDirDlg  lcl.ISaveDialog
	savePathInp lcl.IMemo
	saveFileBuf *bufio.Writer

	fileWindow []*FileWindow
}

var MainWindow TMainWindow

func (m *TMainWindow) FormCreate(sender lcl.IObject) {
	m.SetCaption("ENERGY - XTA Chat UI")
	m.SetPosition(types.PoScreenCenter)
	m.SetWidth(1024)
	m.SetHeight(768)

	m.Constraints().SetMinWidth(types.TConstraintSize(m.Width()))
	m.Constraints().SetMinHeight(types.TConstraintSize(m.Height()))

	png := lcl.NewPngImage()
	png.LoadFromFSFile("assets/icon.png")
	lcl.Application.Icon().Assign(png)
	png.Free()

	m.initMainBox()
}

func (m *TMainWindow) initMainBox() {
	go m.initXTASDK()

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

	// 选择文件
	m.selDirDlg = lcl.NewOpenDialog(m)
	m.selDirDlg.SetOptions(m.selDirDlg.Options().Include(types.OfShowHelp, types.OfAllowMultiSelect))
	m.selDirDlg.SetTitle("XTA - AI SDK 打开文件 多选")

	m.selFileBtn = lcl.NewButton(m)
	m.selFileBtn.SetParent(m)
	m.selFileBtn.SetTop(m.chat.Top() + m.chat.Height() + 3)
	m.selFileBtn.SetWidth(150)
	m.selFileBtn.SetHeight(40)
	m.selFileBtn.SetCaption("选择文件/多选")
	m.selFileBtn.Font().SetSize(12)
	m.selFileBtn.SetLeft(150)
	m.selFileBtn.SetOnClick(m.selectFileOrDir)
	m.selFileBtn.SetAnchors(types.NewSet(types.AkLeft, types.AkBottom))

	// 保存消息
	m.saveDirDlg = lcl.NewSaveDialog(m)
	m.saveDirDlg.SetFilter("文本文件(*.txt)|*.txt|所有文件(*.*)|*.*")
	m.saveDirDlg.SetTitle("XTA - AI SDK 消息保存")
	m.saveChatBtn = lcl.NewButton(m)
	m.saveChatBtn.SetParent(m)
	m.saveChatBtn.SetTop(m.chat.Top() + m.chat.Height() + 3)
	m.saveChatBtn.SetLeft(m.selFileBtn.Left() + m.selFileBtn.Width() + 3)
	m.saveChatBtn.SetCaption("保存消息")
	m.saveChatBtn.SetWidth(150)
	m.saveChatBtn.SetHeight(40)
	m.saveChatBtn.Font().SetSize(12)
	m.saveChatBtn.SetAnchors(types.NewSet(types.AkLeft, types.AkBottom))

	m.saveChatBtn.SetOnClick(func(sender lcl.IObject) {
		if m.saveDirDlg.Execute() {
			m.savePathInp.SetText(m.saveDirDlg.FileName())
		}
	})

	m.savePathInp = lcl.NewMemo(m)
	m.savePathInp.SetParent(m)
	m.savePathInp.SetTop(m.chat.Top() + m.chat.Height() + 3)
	m.savePathInp.SetLeft(m.saveChatBtn.Left() + m.saveChatBtn.Width() + 3)
	m.savePathInp.SetHeight(40)
	m.savePathInp.SetWidth(300)
	m.savePathInp.Font().SetSize(15)
	m.savePathInp.SetWordWrap(false)
	m.savePathInp.SetAnchors(types.NewSet(types.AkLeft, types.AkBottom))
	var savefile *os.File
	m.savePathInp.SetOnChange(func(sender lcl.IObject) {
		if savefile != nil {
			savefile.Close()
			savefile = nil
			m.saveFileBuf = nil
		}
		path := m.savePathInp.Text()
		fe, err := os.Open(path)
		if err == nil {
			defer fe.Close()
			st, err := fe.Stat()
			if err == nil {
				if !st.IsDir() {
					savefile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err == nil {
						m.saveFileBuf = bufio.NewWriter(savefile)
					}
				}
			}
		} else if errors.Is(err, fs.ErrNotExist) {
			savefile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				m.saveFileBuf = bufio.NewWriter(savefile)
			}
		}
	})

	// 窗口显示
	m.SetOnShow(func(sender lcl.IObject) {
		m.chat.SetFocus()
	})

	clearHistory := lcl.NewButton(m)
	clearHistory.SetParent(m)
	clearHistory.SetOnClick(func(sender lcl.IObject) {
		m.ai.History()
	})
}

// 主窗口左侧创建文件项
func (m *TMainWindow) createFileItem(filewindow *FileWindow) {
	btn := lcl.NewButton(m)
	btn.SetParent(m)
	caption := filewindow.fileDesc
	if caption == "" {
		caption = filewindow.filenames
	}
	btn.SetCaption(caption)
	btn.SetOnClick(func(sender lcl.IObject) {
		m.removeFileBtn(filewindow.id)
	})
	btn.SetHint("点击删除该文件项")
	filewindow.fileBtn = btn
	m.fileWindow = append(m.fileWindow, filewindow)
	m.resortFileBtns()
}

func (m *TMainWindow) removeFileBtn(id string) {
	var newwindows []*FileWindow
	for _, fw := range m.fileWindow {
		if fw.id == id {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				fw.fileBtn.Free()
				fw.Close()
			})
		} else {
			newwindows = append(newwindows, fw)
		}
	}
	m.fileWindow = newwindows
	m.resortFileBtns()
}

func (m *TMainWindow) resortFileBtns() {
	for i, fw := range m.fileWindow {
		fw.fileBtn.SetLeft(5)
		fw.fileBtn.SetWidth(130)
		fw.fileBtn.SetHeight(30)
		fw.fileBtn.SetTop(int32(i*30) + 5)
	}
}

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
	})
}

func (m *TMainWindow) selectFileOrDir(sender lcl.IObject) {
	if m.selDirDlg.Execute() {
		files := m.selDirDlg.Files()
		var file []string
		for i := 0; i < int(files.Count()); i++ {
			fmt.Println(files.ValueFromIndex(int32(i)))
			file = append(file, files.ValueFromIndex(int32(i)))
		}
		win := createWindow(file, func(window *FileWindow) {
			m.createFileItem(window)
		})
		win.Show()
	}
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

func nowDatetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
