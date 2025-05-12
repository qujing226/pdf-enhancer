package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DeepSeekConfig 配置DeepSeek API客户端
type DeepSeekConfig struct {
	APIKey      string
	BaseURL     string
	ModelName   string
	MaxTokens   int
	Temperature float64
}

// DeepSeekClient DeepSeek API客户端
type DeepSeekClient struct {
	config     DeepSeekConfig
	httpClient *http.Client
}

// DeepSeekRequest 请求结构
type DeepSeekRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekResponse 响应结构
type DeepSeekResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewDeepSeekClient 创建新的DeepSeek客户端
func NewDeepSeekClient(config DeepSeekConfig) *DeepSeekClient {
	return &DeepSeekClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateSummary 生成报告摘要
func (c *DeepSeekClient) GenerateSummary(ctx context.Context, reportContent string) (string, error) {
	// 构建提示词
	prompt := fmt.Sprintf("请为以下报告生成一个简洁的摘要（不超过200字）:\n\n%s", reportContent)

	// 构建请求
	request := DeepSeekRequest{
		Model: c.config.ModelName,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.config.BaseURL+"/v1/chat/completions",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var response DeepSeekResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查响应是否有效
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("API返回的响应没有内容")
	}

	// 返回摘要内容
	return response.Choices[0].Message.Content, nil
}

// MockGenerateSummary 模拟生成报告摘要（当无法访问DeepSeek API时使用）
func (c *DeepSeekClient) MockGenerateSummary(_ context.Context, reportContent string) (string, error) {
	// 简单地返回一个固定的摘要
	return fmt.Sprintf("这是一份关于%s的报告摘要。该报告包含了重要的财务数据和投资建议。"+
		"根据报告内容，建议投资者关注市场波动并调整投资策略。", reportContent[:20]), nil
}
