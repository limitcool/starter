package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Option 客户端选项
type Option func(*options)

// options 客户端选项
type options struct {
	headers map[string]string
	timeout time.Duration
}

// WithHeader 设置请求头
func WithHeader(key, value string) Option {
	return func(o *options) {
		if o.headers == nil {
			o.headers = make(map[string]string)
		}
		o.headers[key] = value
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// GetJSON 发送GET请求并返回JSON响应
func GetJSON(ctx context.Context, url string, opts ...Option) ([]byte, error) {
	return doRequest(ctx, http.MethodGet, url, nil, opts...)
}

// PostJSON 发送POST请求并返回JSON响应
func PostJSON(ctx context.Context, url string, body []byte, opts ...Option) ([]byte, error) {
	// 添加JSON内容类型头
	opts = append(opts, WithHeader("Content-Type", "application/json"))
	return doRequest(ctx, http.MethodPost, url, bytes.NewReader(body), opts...)
}

// PostForm 发送POST表单请求并返回JSON响应
func PostForm(ctx context.Context, url string, form url.Values, opts ...Option) ([]byte, error) {
	// 添加表单内容类型头
	opts = append(opts, WithHeader("Content-Type", "application/x-www-form-urlencoded"))
	return doRequest(ctx, http.MethodPost, url, strings.NewReader(form.Encode()), opts...)
}

// doRequest 执行HTTP请求
func doRequest(ctx context.Context, method, url string, body io.Reader, opts ...Option) ([]byte, error) {
	// 默认选项
	o := &options{
		timeout: 30 * time.Second, // 默认30秒超时
	}

	// 应用选项
	for _, opt := range opts {
		opt(o)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// 添加请求头
	for k, v := range o.headers {
		req.Header.Set(k, v)
	}

	// 创建客户端
	client := &http.Client{
		Timeout: o.timeout,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	return io.ReadAll(resp.Body)
}
