# XTA AI Go SDK 使用文档

## 项目概述

XTA AI Go SDK 是一个基于 Go 语言开发的库，用于与 AI 服务进行交互。该库支持多种 AI 模型，提供流式和非流式消息处理功能，适用于开发智能对话应用。

## 安装与配置

### 安装

使用 Go 模块安装：

```bash
go get github.com/energye/xta
```

### 配置

在使用 SDK 前，请确保设置 AI API 密钥：

```go
os.Setenv(chat.ENV_AI_API_KEY, "your_api_key")
```

## 快速上手

以下是一个简单的示例，展示了如何使用该库进行流式对话：

```go
package main

import (
    "fmt"
    "os"
    "github.com/energye/xta/chat"
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
                for _, choice := range message.Choices {
                    fmt.Print(choice.Delta.Content)
                }
            }
        }
    })

    ai.SetOnFail(func(message *chat.TResponseError) {
        fmt.Println("fail:", message.Code, message.Message, message.Type)
    })

    ai.ChatStream("你好")
}
```

## 核心功能

### 消息处理

支持以下消息处理方式：

- `Chat(string)`: 发送非流式消息
- `ChatStream(string)`: 发送流式消息
- `ChatRole(string, Role)`: 发送指定角色的消息
- `ChatStreamRole(string, Role)`: 发送指定角色的流式消息

### 模型参数配置

支持配置以下参数：

- `MaxTokens`: 最大生成长度
- `Temperature`: 温度
- `TopP`: Top p
- `Stream`: 是否启用流式输出
- `Stop`: 停止词列表
- `Seed`: 随机种子
- 等等

### 回调函数

提供以下回调函数用于处理响应：

- `SetOnReceive(func(*TResponse))`: 设置接收消息的回调
- `SetOnFail(func(*TResponseError))`: 设置处理错误的回调

## 高级功能

### 模型切换

支持动态切换模型：

```go
ai.SetModel(chat.GITEE_AI_DEEPSEEK_R1_DISTILL_QWEN_32B)
```

### 工具支持

支持集成工具函数，扩展 AI 的功能。

## 示例代码

以下是一些示例代码，展示了如何使用该 SDK：

### 流式对话

```go
ai.ChatStream("你好")
```

### 非流式对话

```go
ai.Chat("你好")
```

### 配置模型参数

```go
ai.MetaData().MaxTokens = 100
ai.MetaData().Temperature = 0.7
```

### 设置停止词

```go
ai.MetaData().Stop = []string{"stop", "end"}
```

## 常见问题

### 问题：API 密钥未设置

**解决方法**：

确保在运行前设置 API 密钥：

```go
os.Setenv(chat.ENV_AI_API_KEY, "your_api_key")
```

### 问题：模型未支持

**解决方法**：

检查支持的模型列表，确保使用的模型名称正确。

### 问题：网络请求失败

**解决方法**：

检查网络连接，确保 API 地址正确，并处理可能的网络错误。

## 总结

XTA AI Go SDK 提供了一个方便的接口，用于与 XTA AI 服务交互。通过配置不同的模型和参数，开发者可以灵活地构建各种智能对话应用。如需进一步帮助，请参考项目文档或联系支持团队。