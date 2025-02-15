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
	"log"
	"net/http"
	"net/url"
	"os"
)

// 模型配置

// 服务地址
const GITEE_AI = "https://ai.gitee.com/v1"

// 服务接口
const GITEE_AI_API = "/chat/completions"

// gitee ai 模型列表
type GiteeAIModelNameEnum string

const (
	GITEE_AI_YI_34B_CHAT                   GiteeAIModelNameEnum = "Yi-34B-Chat"
	GITEE_AI_INTERNVL2_8B                  GiteeAIModelNameEnum = "InternVL2-8B"
	GITEE_AI_INTERNVL2_5_78B               GiteeAIModelNameEnum = "InternVL2.5-78B"
	GITEE_AI_DEEPSEEK_CODER_33B_INSTRUCT   GiteeAIModelNameEnum = "deepseek-coder-33B-instruct"
	GITEE_AI_INTERNVL2_5_26B               GiteeAIModelNameEnum = "InternVL2.5-26B"
	GITEE_AI_DEEPSEEK_R1_DISTILL_QWEN_1_5B GiteeAIModelNameEnum = "DeepSeek-R1-Distill-Qwen-1.5B"
	GITEE_AI_QWEN2_VL_72B                  GiteeAIModelNameEnum = "Qwen2-VL-72B"
	GITEE_AI_QWEN2_5_32B_INSTRUCT          GiteeAIModelNameEnum = "Qwen2.5-32B-Instruct"
	GITEE_AI_GLM_4_9B_CHAT                 GiteeAIModelNameEnum = "glm-4-9b-chat"
	GITEE_AI_QWQ_32B_PREVIEW               GiteeAIModelNameEnum = "QwQ-32B-Preview"
	GITEE_AI_CODEGEEX4_ALL_9B              GiteeAIModelNameEnum = "codegeex4-all-9b"
	GITEE_AI_QWEN2_5_CODER_32B_INSTRUCT    GiteeAIModelNameEnum = "Qwen2.5-Coder-32B-Instruct"
	GITEE_AI_DEEPSEEK_R1                   GiteeAIModelNameEnum = "DeepSeek-R1"
	GITEE_AI_QWEN2_5_72B_INSTRUCT          GiteeAIModelNameEnum = "Qwen2.5-72B-Instruct"
	GITEE_AI_QWEN2_5_7B_INSTRUCT           GiteeAIModelNameEnum = "Qwen2.5-7B-Instruct"
	GITEE_AI_DEEPSEEK_R1_DISTILL_QWEN_7B   GiteeAIModelNameEnum = "DeepSeek-R1-Distill-Qwen-7B"
	GITEE_AI_QWEN2_5_CODER_14B_INSTRUCT    GiteeAIModelNameEnum = "Qwen2.5-Coder-14B-Instruct"
	GITEE_AI_DEEPSEEK_R1_DISTILL_QWEN_32B  GiteeAIModelNameEnum = "DeepSeek-R1-Distill-Qwen-32B"
	GITEE_AI_QWEN2_72B_INSTRUCT            GiteeAIModelNameEnum = "Qwen2-72B-Instruct"
	GITEE_AI_CODE_RACCOON_V1               GiteeAIModelNameEnum = "code-raccoon-v1"
	GITEE_AI_QWEN2_7B_INSTRUCT             GiteeAIModelNameEnum = "Qwen2-7B-Instruct"
	GITEE_AI_DEEPSEEK_V3                   GiteeAIModelNameEnum = "DeepSeek-V3"
	GITEE_AI_QWEN2_5_14B_INSTRUCT          GiteeAIModelNameEnum = "Qwen2.5-14B-Instruct"
	GITEE_AI_DEEPSEEK_R1_DISTILL_QWEN_14B  GiteeAIModelNameEnum = "DeepSeek-R1-Distill-Qwen-14B"
)

func GiteeAIModels() []string {
	return []string{
		string(GITEE_AI_YI_34B_CHAT),
		string(GITEE_AI_INTERNVL2_8B),
		string(GITEE_AI_INTERNVL2_5_78B),
		string(GITEE_AI_DEEPSEEK_CODER_33B_INSTRUCT),
		string(GITEE_AI_INTERNVL2_5_26B),
		string(GITEE_AI_DEEPSEEK_R1_DISTILL_QWEN_1_5B),
		string(GITEE_AI_QWEN2_VL_72B),
		string(GITEE_AI_QWEN2_5_32B_INSTRUCT),
		string(GITEE_AI_GLM_4_9B_CHAT),
		string(GITEE_AI_QWQ_32B_PREVIEW),
		string(GITEE_AI_CODEGEEX4_ALL_9B),
		string(GITEE_AI_QWEN2_5_CODER_32B_INSTRUCT),
		string(GITEE_AI_DEEPSEEK_R1),
		string(GITEE_AI_QWEN2_5_72B_INSTRUCT),
		string(GITEE_AI_QWEN2_5_7B_INSTRUCT),
		string(GITEE_AI_DEEPSEEK_R1_DISTILL_QWEN_7B),
		string(GITEE_AI_QWEN2_5_CODER_14B_INSTRUCT),
		string(GITEE_AI_DEEPSEEK_R1_DISTILL_QWEN_32B),
		string(GITEE_AI_QWEN2_72B_INSTRUCT),
		string(GITEE_AI_CODE_RACCOON_V1),
		string(GITEE_AI_QWEN2_7B_INSTRUCT),
		string(GITEE_AI_DEEPSEEK_V3),
		string(GITEE_AI_QWEN2_5_14B_INSTRUCT),
		string(GITEE_AI_DEEPSEEK_R1_DISTILL_QWEN_14B),
	}
}

// IGiteeAI gitee ai 模型接口
type IGiteeAI interface {
	IAI
	SetModel(name GiteeAIModelNameEnum)       // 设置当前模型
	Model() GiteeAIModelNameEnum              // 返回当前模型
	ChatRole(content string, role Role)       // 带有角色的聊天, 发送消息并以普通方式全量返回
	ChatStreamRole(content string, role Role) // 带有角色的聊天, 发送消息并以流方式返回
	Chat(content string)                      // 发送消息并以普通方式全量返回
	ChatStream(content string)                // 发送消息并以流方式返回
	Options() *Options                        // 返回当前 AI 选项
}

// Value 返回模型枚举值
func (m GiteeAIModelNameEnum) Value() string {
	return string(m)
}

// DefaultGiteeAIMetaData 默认模型参数配置
var DefaultGiteeAIMetaData = MetaData{
	Model:    GITEE_AI_DEEPSEEK_R1_DISTILL_QWEN_32B.Value(),
	Messages: make(Messages, 0),
}

// DefaultGiteeAIOptions 默认选项配置
var DefaultGiteeAIOptions = Options{
	BaseURL: GITEE_AI,
	API:     GITEE_AI_API,
	APIKey:  os.Getenv(ENV_AI_API_KEY),
}

// GiteeAI Gitee AI 实现
type GiteeAI struct {
	AIBase
	options       *Options // AI 选项
	isSupportTool bool     // 是否支持工具
	metaData      MetaData // 元数据参数
	header        http.Header
}

// NewGiteeAI 创建一个 Gitee AI
func NewGiteeAI(options Options, isSupportTool bool) IGiteeAI {
	ai := &GiteeAI{
		options:       &options,
		isSupportTool: isSupportTool,
		metaData:      DefaultGiteeAIMetaData,
		header:        make(http.Header),
	}
	return ai
}

func (m *GiteeAI) API() string {
	result, err := url.JoinPath(m.options.BaseURL, m.options.API)
	if err != nil {
		log.Println(err)
	}
	return result
}

func (m *GiteeAI) Name() string {
	return m.MetaData().Model
}

func (m *GiteeAI) IsSupportTool() bool {
	return m.isSupportTool
}

func (m *GiteeAI) APIKey() string {
	return m.options.APIKey
}

func (m *GiteeAI) Options() *Options {
	return m.options
}

func (m *GiteeAI) MetaData() *MetaData {
	return &m.metaData
}

func (m *GiteeAI) SetModel(name GiteeAIModelNameEnum) {
	m.MetaData().Model = string(name)
}

func (m *GiteeAI) Model() GiteeAIModelNameEnum {
	return GiteeAIModelNameEnum(m.MetaData().Model)
}

func (m *GiteeAI) ChatRole(content string, role Role) {
	m.MetaData().Stream = false
	m.history.Add(Message{Role: role, Content: content})
	m.Request()
}

func (m *GiteeAI) ChatStreamRole(content string, role Role) {
	m.MetaData().Stream = true
	m.history.Add(Message{Role: role, Content: content})
	m.Request()
}

func (m *GiteeAI) Chat(content string) {
	m.ChatRole(content, RoleUser)
}

func (m *GiteeAI) ChatStream(content string) {
	m.ChatStreamRole(content, RoleUser)
}

func (m *GiteeAI) Header() http.Header {
	return m.header
}

func (m *GiteeAI) Request() {
	m.Header().Set("Authorization", "Bearer "+m.APIKey())
	m.Header().Set("Content-Type", "application/json")
	HttpPost(m)
}
