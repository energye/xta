# CEF Go 库使用文档

## 概述

CEF（Chromium Embedded Framework）是一个用于在应用中嵌入Chromium浏览器的框架。本库是一个基于Go语言的CEF实现，旨在为开发者提供简便的方式在Go应用中集成网页浏览功能。通过本库，开发者可以轻松创建支持HTML5、JavaScript和现代Web技术的桌面应用。

## 安装与配置

### 安装依赖

在使用本库之前，请确保已安装以下依赖项：

1. **golcl**: 用于与LCL（ Lazarus Component Library）交互。
2. **energy**: 提供底层CEF实现和工具函数。
3. 其他系统依赖：如CEF二进制文件和相关依赖库。

#### 使用Go模块安装

```bash
go get github.com/energye/golcl/...
go get github.com/energye/energy/v2/...
```

### 配置CEF环境

在代码中全局初始化CEF环境：

```go
import (
    "github.com/energye/energy/v2/cef"
    "github.com/energye/energy/v2/common"
    _ "github.com/energye/energy/v2/examples/syso"
)

func main() {
    // 全局初始化，每个应用必须调用
    cef.GlobalInit(nil, resources)
    // 创建CEF应用实例
    cefApp := cef.NewApplication()
    // ... 其他配置
}
```

## 快速上手

### 示例代码

以下是一个简单的示例，展示了如何使用本库创建一个浏览器应用：

```go
package main

import (
    "embed"
    "fmt"
    "github.com/energye/energy/v2/cef"
    "github.com/energye/energy/v2/common"
    _ "github.com/energye/energy/v2/examples/syso"
    "github.com/energye/energy/v2/pkgs/assetserve"
    "github.com/energye/golcl/lcl/api"
)

//go:embed resources
var resources embed.FS

func main() {
    // 全局初始化，每个应用必须调用
    cef.GlobalInit(nil, resources)
    // 创建CEF应用实例
    cefApp := cef.NewApplication()

    // 设置浏览器窗口的URL
    cef.BrowserWindow.Config.Url = "http://localhost:22022/audio-video.html"

    // 根据平台设置图标
    if common.IsLinux() && api.WidgetUI().IsGTK3() {
        cef.BrowserWindow.Config.IconFS = "resources/icon.png"
    } else {
        cef.BrowserWindow.Config.IconFS = "resources/icon.ico"
    }

    // 主进程启动成功后的回调
    cef.SetBrowserProcessStartAfterCallback(func(b bool) {
        fmt.Println("主进程启动 创建一个内置HTTP服务")
        // 启动内置HTTP服务器
        server := assetserve.NewAssetsHttpServer()
        server.PORT = 22022
        server.AssetsFSName = "resources"
        server.Assets = resources
        go server.StartHttpServer()
    })

    // 运行CEF应用
    cef.Run(cefApp)
}
```

### 代码解析

1. **全局初始化**：`cef.GlobalInit(nil, resources)` 初始化CEF环境，并传入嵌入的资源文件系统。
2. **创建应用实例**：`cef.NewApplication()` 创建一个CEF应用实例。
3. **设置URL**：`cef.BrowserWindow.Config.Url` 设置浏览器窗口加载的初始URL。
4. **设置图标**：根据平台设置不同的图标资源路径。
5. **启动HTTP服务器**：在主进程启动后，创建并启动一个内置HTTP服务器，用于服务嵌入的静态资源。
6. **运行应用**：`cef.Run(cefApp)` 启动CEF应用，进入消息循环。

## 核心功能

### 创建和管理CEF应用

#### 创建应用实例

```go
cefApp := cef.NewApplication()
```

- `NewApplication` 创建一个新的CEF应用实例。默认情况下，会注册一些默认事件处理函数，并初始化一些基本设置。

#### 启动应用

```go
cef.Run(cefApp)
```

- `Run` 函数负责启动CEF应用，初始化窗口组件，并进入消息循环。在Linux平台上，默认使用Views Framework窗口组件；在Windows和macOS上，默认使用LCL窗口组件。

### 配置CEF应用

#### 设置缓存路径

```go
cefApp.SetCache("path/to/cache")
```

- `SetCache` 方法设置浏览器的缓存目录。

#### 设置用户代理

```go
cefApp.SetUserAgent("MyCustomUserAgent")
```

- `SetUserAgent` 方法设置浏览器发送的HTTP用户代理字符串。

#### 启用/禁用GPU加速

```go
cefApp.SetEnableGPU(true)
```

- `SetEnableGPU` 方法启用或禁用GPU加速渲染。

#### 设置语言

```go
cefApp.SetLocale(cef.LANGUAGE_zh_CN)
```

- `SetLocale` 方法设置浏览器的显示语言。

### 控制消息循环

#### 运行消息循环

```go
cefApp.RunMessageLoop()
```

- `RunMessageLoop` 方法进入消息循环，处理窗口事件和用户输入。在LCL窗口组件中，消息循环由LCL管理，在Views Framework中则由CEF管理。

#### 退出消息循环

```go
cefApp.QuitMessageLoop()
```

- `QuitMessageLoop` 方法退出消息循环，通常在应用关闭时调用。

### 处理事件

#### 注册自定义协议

```go
cefApp.SetOnRegCustomSchemes(func(registrar *cef.TCefSchemeRegistrarRef) {
    registrar.AddCustomScheme("mycustomscheme", cef.CEF_SCHEME_OPTION_STANDARD)
})
```

- `SetOnRegCustomSchemes` 方法注册自定义URL协议，允许应用处理特定格式的URL。

#### 处理上下文创建事件

```go
cefApp.SetOnContextCreated(func(browser *cef.ICefBrowser, frame *cef.ICefFrame, context *cef.ICefV8Context) bool {
    fmt.Println("Context created")
    return false
})
```

- `SetOnContextCreated` 方法设置上下文创建事件的回调函数，适用于需要在JavaScript上下文中注入代码的场景。

## 高级主题

### 异步调用

在多线程或多进程环境下，确保在主线程中执行UI相关操作至关重要。本库提供以下方法处理异步调用：

#### 异步调用

```go
cef.QueueAsyncCall(func(id int) {
    fmt.Println("异步任务执行中")
})
```

- `QueueAsyncCall` 方法将任务排队，确保在主线程中异步执行。

#### 同步调用

```go
cef.QueueSyncCall(func(id int) {
    fmt.Println("同步任务执行中")
})
```

- `QueueSyncCall` 方法同步执行任务，阻塞直到任务完成。

### 处理窗口事件

#### 创建窗口

```go
cef.BrowserWindow.createFormAndRun()
```

- 在LCL窗口组件中，调用`createFormAndRun`方法创建并显示主窗口。

#### 处理窗口关闭事件

```go
cef.BrowserWindow.Config.OnClose = func() {
    fmt.Println("窗口关闭")
    cef.BrowserWindow.Destroy()
}
```

- 设置窗口关闭事件的回调函数，释放资源并清理状态。

### 自定义CEF行为

#### 禁用安全浏览

```go
cefApp.SetDisableSafeBrowsing(true)
```

- 禁用Chrome的安全浏览功能。

#### 启用开发者工具

```go
cefApp.SetRemoteDebuggingPort(9222)
```

- 启用远程调试功能，允许通过Chrome开发者工具调试应用。

## 常见问题

### 问题：CEF文件未找到

**解决方法**：

1. 确保CEF二进制文件存在于正确的路径下。
2. 在代码中设置CEF框架路径：

```go
cefApp.SetFrameworkDirPath("/path/to/cef/framework")
```

### 问题：窗口无法显示

**解决方法**：

1. 确保所有依赖库已正确安装。
2. 在Linux平台上，检查是否安装了必要的GTK库：

```bash
sudo apt-get install libgtk-3-dev
```

3. 确保CEF版本与操作系统兼容。

### 问题：JavaScript无法执行

**解决方法**：

1. 确保JavaScript未被禁用：

```go
cefApp.SetDisableJavascript(false)
```

2. 检查控制台日志，查看是否有JavaScript错误。

### 问题：HTTP服务器无法启动

**解决方法**：

1. 确保端口未被占用：

```go
server.PORT = 22022 // 尝试更换端口
```

2. 检查防火墙设置，确保端口开放。

## 总结

通过本库，开发者可以轻松在Go应用中集成CEF，创建功能丰富的桌面浏览器应用。文档详细介绍了从安装配置到高级功能的各个方面，帮助开发者快速上手并解决常见问题。如需进一步帮助，请参考CEF官方文档或社区资源。