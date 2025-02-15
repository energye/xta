package main

import (
	"embed"
	"fmt"
	"github.com/energye/lcl/inits"
	"github.com/energye/lcl/lcl"
	"xta/xtaui/window"
)

//go:embed assets
var resources embed.FS

//go:embed libs
var lib embed.FS

func main() {
	inits.Init(lib, resources)
	lcl.Application.Initialize()
	lcl.Application.SetOnException(func(sender lcl.IObject, e lcl.IException) {
		fmt.Println("Exception:", e.ToString())
	})
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.CreateForm(&window.MainWindow)
	lcl.Application.Run()
}
