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

import (
	"encoding/json"
	"log"
)

var Debug bool

const (
	AUTHORIZATION = "Authorization"
	CONTENT_TYPE  = "Content-Type"
)

const ENV_AI_API_KEY = "AI_API_KEY"

// Role 角色
type Role string

const (
	RoleUser      = "user"      //
	RoleSystem    = "system"    //
	RoleAssistant = "assistant" //
)

// FuncParameterType 工具函数参数类型
type FuncParameterType string

const (
	FptInteger FuncParameterType = "integer"
	FptString  FuncParameterType = "string"
	FptFloat   FuncParameterType = "float"
)

// Messages 消息列表
type Messages []Message

func (m *Messages) Add(message Message) {
	*m = append(*m, message)
}

// ResponseFormat 响应格式
type ResponseFormat map[string]string

func (m *ResponseFormat) Add(name, value string) {
	if *m == nil {
		*m = make(map[string]string)
	}
	(*m)[name] = value
}

type LogitBias map[string]float32

// Add 值范围 -100 ~ 100
func (m *LogitBias) Add(name string, value float32) {
	if *m == nil {
		*m = make(map[string]float32)
	}
	if value >= -100 && value <= 100 {
		(*m)[name] = value
	}
}

// FuncParameters 模型调用的工具函数参数列表
type FuncParameters map[string]FuncParameter

func (m *FuncParameters) Add(name string, type_ FuncParameterType, description string) {
	if *m == nil {
		*m = make(map[string]FuncParameter)
	}
	(*m)[name] = FuncParameter{Type: type_, Description: description}
}

func (m *FuncParameters) AddString(name string, description string) {
	m.Add(name, FptString, description)
}

func (m *FuncParameters) AddInteger(name string, description string) {
	m.Add(name, FptInteger, description)
}

func (m *FuncParameters) AddFloat(name string, description string) {
	m.Add(name, FptFloat, description)
}

// MetaData 元数据参数配置
type MetaData struct {
	Model            string         `json:"model,omitempty"`             // 模型 required
	Messages         Messages       `json:"messages,omitempty"`          // 包含当前会话所有消息的列表
	Stream           bool           `json:"stream"`                      // 是否是流式输出
	MaxTokens        int            `json:"max_tokens,omitempty"`        // 最大生成长度 default: 0
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"` // 频率惩罚 default: 0
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`  // 存在惩罚 default: 0
	Stop             []string       `json:"stop,omitempty"`              // 停止词
	Temperature      float32        `json:"temperature,omitempty"`       // 温度 default: 1
	TopP             float32        `json:"top_p,omitempty"`             // Top p default: 1
	TopLogprobs      int            `json:"top_logprobs,omitempty"`      // Top Logprobs default: 0
	ResponseFormat   ResponseFormat `json:"response_format,omitempty"`   // 响应格式
	Seed             int            `json:"seed,omitempty"`              // 随机种子 default: 0
	N                int            `json:"n,omitempty"`                 // 生成数量 default: 1
	LogitBias        LogitBias      `json:"logit_bias,omitempty"`        // 对数偏差 additional properties double min: -100 max: 100 default: {"1000": 99.9, "1001": -99.9}
	User             string         `json:"user,omitempty"`              // 客户侧的用户标识 string | nullable
	Tools            Tools          `json:"tools,omitempty"`             // 模型可以调用的工具列表。目前仅支持函数作为工具 array object[]
	GuidedJSON       any            `json:"guided_json,omitempty"`       // 如果指定，输出将遵循JSON模式 anyof
	GuidedChoice     []string       `json:"guided_choice,omitempty"`     // 如果指定，输出将恰好是其中一个选项
}

// Message 消息
type Message struct {
	Role    Role   `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// Tools 工具列表
type Tools []*Tool

func (m *Tools) Add(tool *Tool) {
	*m = append(*m, tool)
}

// Tool 模型调用的工具
type Tool struct {
	Type     string    `json:"type,omitempty"`     // const: function default: function
	Function *Function `json:"function,omitempty"` // 包含函数的详细信息 required
}

func NewTool() *Tool {
	return &Tool{Type: "function"}
}

func (m *Tool) AddFunc(name, desc string) *Function {
	m.Function = &Function{Name: name, Description: desc}
	return m.Function
}

// Function 模型调用的工具函数
type Function struct {
	Name        string         `json:"name,omitempty"`        // 函数的名称 required
	Description string         `json:"description,omitempty"` // 函数的描述
	Parameters  FuncParameters `json:"parameters,omitempty"`  // 函数接受的参数
}

// FuncParameter 模型调用的工具函数参数
type FuncParameter struct {
	Type        FuncParameterType `json:"type,omitempty"`
	Description string            `json:"description,omitempty"`
}

// ToJSON 模型元数据配置转成 JSON
func (m *MetaData) ToJSON() (data []byte) {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	return
}
