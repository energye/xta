package window

import (
	"bytes"
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"io/ioutil"
	"path/filepath"
	"time"
)

type FileWindow struct {
	lcl.IForm
	fileBtn     lcl.IButton
	text        lcl.IMemo
	filenames   string
	id          string
	fileDesc    string
	fileContent []string
}

func createWindow(files []string, ok func(window *FileWindow)) *FileWindow {
	form := lcl.NewForm(lcl.Application)
	form.SetPosition(types.PoScreenCenter)
	form.SetWidth(300)
	form.SetHeight(200)
	form.SetBorderStyleForFormBorderStyle(types.BsNone)
	form.SetColor(colors.ClAzure)

	window := &FileWindow{IForm: form, id: time.Now().String()}

	fileLabel := lcl.NewLabel(form)
	fileLabel.SetParent(form)
	fileLabel.SetWidth(form.Width())
	fileLabel.SetLeft(5)
	fileLabel.SetTop(5)
	fileLabel.SetCaption("文件描述和作用")

	window.text = lcl.NewMemo(form)
	window.text.SetParent(form)
	window.text.SetBorderStyle(types.BsSingle)
	window.text.SetHeight(170)
	window.text.SetWidth(300)
	window.text.SetTop(30)
	window.text.Font().SetSize(12)

	okBtn := lcl.NewButton(form)
	okBtn.SetParent(form)
	okBtn.SetCaption("确认")
	okBtn.SetWidth(50)
	okBtn.SetTop(form.Height() - 30)
	okBtn.SetLeft(form.Width() - 60)
	okBtn.SetOnClick(func(sender lcl.IObject) {
		// 文本描述 + 读取文件内容
		for i, file := range files {
			_, name := filepath.Split(file)
			if i > 0 {
				window.filenames += ", "
			}
			window.filenames += name
			data, err := ioutil.ReadFile(file)
			if err == nil {
				buf := bytes.Buffer{}
				buf.WriteString(name + "\n")
				buf.Write(data)
				buf.WriteString("\n")
				fmt.Println("文件:", file, "大小:", buf.Len())
				window.fileContent = append(window.fileContent, buf.String())
			}
		}
		window.fileDesc = window.text.Text()
		form.Hide()
		if ok != nil {
			ok(window)
		}
	})

	cancelBtn := lcl.NewButton(form)
	cancelBtn.SetParent(form)
	cancelBtn.SetCaption("取消")
	cancelBtn.SetWidth(50)
	cancelBtn.SetTop(form.Height() - 30)
	cancelBtn.SetLeft(form.Width() - 120)
	cancelBtn.SetOnClick(func(sender lcl.IObject) {
		form.Close()
	})
	return window
}
