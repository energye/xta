//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package chat

import "net/http"

// IAI 模型接口
type IAI interface {
	API() string
	Name() string
	APIKey() string
	IsSupportTool() bool
	OnFail() TOnFail
	OnReceive() TOnReceive
	SetOnFail(fn TOnFail)
	SetOnReceive(fn TOnReceive)
	MetaData() *MetaData
	Header() http.Header
	History() *Messages
}

// IGiteeAI gitee ai 模型接口
type IGiteeAI interface {
	IAI
	SetModel(name GiteeAIModelNameEnum)
	Model() GiteeAIModelNameEnum
	ChatRole(content string, role Role)
	ChatStreamRole(content string, role Role)
	Chat(content string)
	ChatStream(content string)
}

// AIBase 基类
type AIBase struct {
	fail    TOnFail    // 回调函数, 会话返回结果
	receive TOnReceive // 回调函数, 成功完成时调用
	history Messages   // 历史消息
}

func (m *AIBase) OnFail() TOnFail {
	return m.fail
}

func (m *AIBase) OnReceive() TOnReceive {
	return m.receive
}

func (m *AIBase) SetOnFail(fn TOnFail) {
	m.fail = fn
}

func (m *AIBase) SetOnReceive(fn TOnReceive) {
	m.receive = fn
}

func (m *AIBase) History() *Messages {
	return &m.history
}
