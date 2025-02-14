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

import "bytes"

// TOnReceive 接收成功消息事件
type TOnReceive func(message *TResponse)

// TOnFail 接收失败消息事件
type TOnFail func(err *TResponseError)

type Choices []TChoice

func (m Choices) ToStringArray() (result []string) {
	result = make([]string, len(m))
	for i, choice := range m {
		result[i] = choice.Delta.Content
	}
	return
}
func (m Choices) ToString() string {
	buf := bytes.Buffer{}
	for _, choice := range m {
		buf.WriteString(choice.Delta.Content)
	}
	return buf.String()
}

type TResponse struct {
	Id        string  `json:"id"`
	Object    string  `json:"object"`
	Created   int     `json:"created"`
	Model     string  `json:"model"`
	Usage     TUsage  `json:"usage"`
	Choices   Choices `json:"choices"`
	Error     string  `json:"error"`
	ErrorType string  `json:"error_type"`
}

type TUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type TChoice struct {
	FinishReason string            `json:"finish_reason"`
	Index        int               `json:"index"`
	Message      map[string]string `json:"message"`
	Logprobs     any               `json:"logprobs"`
	Delta        TDelta            `json:"delta"`
}

type TDelta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TMessage struct {
	Success bool
	Content []string
}

type TError struct {
	Error TResponseError `json:"error"`
}

type TResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}
