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

// IAI 基础接口
type IAI interface {
	API() string                // AI API 服务地址
	Name() string               // 模型名称
	APIKey() string             // 服务商 API KEY
	IsSupportTool() bool        // 模型是否支持工具
	OnFail() TOnFail            // 消息接收时, 返回失败时回调函数
	OnReceive() TOnReceive      // 消息接收时, 返回成功时消息接收
	SetOnFail(fn TOnFail)       // 设置 失败回调函数
	SetOnReceive(fn TOnReceive) // 设置 消息接收回调函数
	MetaData() *MetaData        // 模型源数据参数
	Header() http.Header        // 请求头
	History() *Messages         // 历史消息列表, 包括当前消息
}

// AIBase 基础实现
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
