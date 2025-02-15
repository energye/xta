package window

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/rtl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/xta/chat"
	"io/fs"
	"os"
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

	title string
}

var MainWindow TMainWindow

func (m *TMainWindow) FormCreate(sender lcl.IObject) {
	m.title = "ENERGY - XTA Chat UI"
	m.SetCaption(m.title)
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

	openURL := lcl.NewLinkLabel(m)
	openURL.SetParent(m)
	openURL.SetCaption(`<a href="https://ai.gitee.com/models">Gitee AI API 获取</a>`)
	openURL.SetAlign(types.AlRight)
	openURL.SetTop(5)
	openURL.Font().SetSize(12)
	openURL.SetOnLinkClick(func(sender lcl.IObject, link string, linktype types.TSysLinkType) {
		rtl.SysOpen(link)
	})

	modules := lcl.NewComboBox(m)
	modules.SetParent(m)
	modules.SetLeft(150)
	modules.Items().AddStrings2(chat.GiteeAIModels())
	modules.SetItemIndex(17)
	modules.SetHeight(35)
	modules.SetWidth(300)
	modules.Font().SetSize(12)
	modules.SetOnChange(func(sender lcl.IObject) {
		module := chat.GiteeAIModelNameEnum(modules.Items().Strings(modules.ItemIndex()))
		m.ai.SetModel(module)
		m.message.Lines().Add("模型: " + m.ai.Name())
		m.SetCaption(m.title + " " + m.ai.Name())
	})

	apiKey := lcl.NewEditButton(m)
	apiKey.SetParent(m)
	apiKey.SetLeft(modules.Left() + modules.Width() + 5)
	apiKey.SetPasswordChar(uint16('*'))
	apiKey.SetHeight(35)
	apiKey.SetWidth(200)
	apiKey.Font().SetSize(12)
	apiKey.Button().SetCaption("API KEY")
	//apiKey.Button().SetLeft(100)
	apiKey.Button().SetWidth(80)
	apiKey.SetOnClick(func(sender lcl.IObject) {
		m.ai.Options().APIKey = apiKey.Text()
	})

	// 消息
	m.message = lcl.NewMemo(m)
	m.message.SetParent(m)
	m.message.SetTop(40)
	m.message.SetLeft(150)
	m.message.SetWidth(m.Width() - 150)
	m.message.SetHeight(m.Height() - 190)
	m.message.SetBorderStyle(types.BsNone)
	m.message.SetReadOnly(true)
	m.message.SetScrollBars(types.SsAutoBoth)
	m.message.SetWordWrap(true)
	m.message.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight, types.AkBottom))

	// 聊天

	m.chat = lcl.NewMemo(m)
	m.chat.SetParent(m)
	m.chat.SetBorderStyle(types.BsNone)
	m.chat.SetScrollBars(types.SsAutoBoth)
	m.chat.SetTop(m.message.Top() + m.message.Height() + 3)
	m.chat.SetLeft(150)
	m.chat.SetWidth(m.Width())
	m.chat.SetHeight(100)
	m.chat.SetAnchors(types.NewSet(types.AkLeft, types.AkRight, types.AkBottom))
	chatLabel := lcl.NewLabel(m)
	chatLabel.SetParent(m)
	chatLabel.SetCaption("发送消息")
	chatLabel.SetTop(m.chat.Top() + 40)
	chatLabel.SetLeft(20)
	chatLabel.Font().SetSize(18)
	chatLabel.Font().SetColor(colors.ClGray)
	chatLabel.SetAnchors(types.NewSet(types.AkLeft, types.AkBottom))

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
	m.saveChatBtn.SetWidth(100)
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
		apiKey.SetText(m.ai.APIKey())
	})

	//clearHistory := lcl.NewButton(m)
	//clearHistory.SetParent(m)
	//clearHistory.SetOnClick(func(sender lcl.IObject) {
	//	m.ai.History()
	//})
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

func nowDatetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
