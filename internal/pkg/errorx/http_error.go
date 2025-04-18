package errorx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTP错误类型常量
const (
	HTTPErrorTypeTimeout      = "timeout"
	HTTPErrorTypeConnection   = "connection"
	HTTPErrorTypeServer       = "server"
	HTTPErrorTypeClient       = "client"
	HTTPErrorTypeUnauthorized = "unauthorized"
	HTTPErrorTypeForbidden    = "forbidden"
	HTTPErrorTypeNotFound     = "not_found"
	HTTPErrorTypeUnknown      = "unknown"
)

// HTTPError HTTP错误结构
type HTTPError struct {
	StatusCode int               // HTTP状态码
	URL        string            // 请求URL
	Method     string            // 请求方法
	Headers    map[string]string // 响应头
	Body       string            // 响应体
	Duration   time.Duration     // 请求耗时
	Err        error             // 原始错误
}

// Error 实现error接口
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %s %s failed with status %d: %s", e.Method, e.URL, e.StatusCode, e.Body)
}

// Unwrap 实现errors.Unwrap接口
func (e *HTTPError) Unwrap() error {
	return e.Err
}

// NewHTTPError 创建HTTP错误
func NewHTTPError(resp *http.Response, err error, duration time.Duration) *HTTPError {
	if err != nil && resp == nil {
		return &HTTPError{
			StatusCode: 0,
			URL:        "",
			Method:     "",
			Headers:    make(map[string]string),
			Body:       err.Error(),
			Duration:   duration,
			Err:        err,
		}
	}

	// 读取响应体
	var body string
	if resp != nil && resp.Body != nil {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr == nil {
			body = string(bodyBytes)
			// 重置响应体，以便后续处理
			resp.Body = io.NopCloser(strings.NewReader(body))
		}
	}

	// 提取响应头
	headers := make(map[string]string)
	if resp != nil && resp.Header != nil {
		for k, v := range resp.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}
	}

	return &HTTPError{
		StatusCode: resp.StatusCode,
		URL:        resp.Request.URL.String(),
		Method:     resp.Request.Method,
		Headers:    headers,
		Body:       body,
		Duration:   duration,
		Err:        err,
	}
}

// HandleHTTPError 处理HTTP错误
func HandleHTTPError(err error) error {
	if err == nil {
		return nil
	}

	// 检查是否为上下文取消或超时
	if err == context.Canceled {
		tmpErr := ErrUnknown.WithMsg("请求被取消")
		return tmpErr.WithError(err)
	}
	if err == context.DeadlineExceeded {
		tmpErr := ErrUnknown.WithMsg("请求超时")
		return tmpErr.WithError(err)
	}

	// 检查是否为HTTPError
	var httpErr *HTTPError
	if he, ok := err.(*HTTPError); ok {
		httpErr = he
	} else {
		// 非HTTPError，直接返回通用错误
		return ErrUnknown.WithError(err)
	}

	// 根据状态码处理不同类型的HTTP错误
	switch {
	case httpErr.StatusCode == 0:
		// 连接错误
		if strings.Contains(httpErr.Body, "timeout") || strings.Contains(httpErr.Body, "deadline") {
			tmpErr := ErrUnknown.WithMsg("请求超时")
			return tmpErr.WithError(err)
		}
		tmpErr := ErrUnknown.WithMsg("连接失败")
		return tmpErr.WithError(err)
	case httpErr.StatusCode == http.StatusUnauthorized:
		tmpErr := ErrUserAuthFailed.WithMsg("认证失败")
		return tmpErr.WithError(err)
	case httpErr.StatusCode == http.StatusForbidden:
		tmpErr := ErrAccessDenied.WithMsg("访问被拒绝")
		return tmpErr.WithError(err)
	case httpErr.StatusCode == http.StatusNotFound:
		tmpErr := ErrNotFound.WithMsg("资源不存在")
		return tmpErr.WithError(err)
	case httpErr.StatusCode >= 500:
		tmpErr := ErrUnknown.WithMsg("服务器错误")
		return tmpErr.WithError(err)
	case httpErr.StatusCode >= 400:
		tmpErr := ErrInvalidParams.WithMsg("请求参数错误")
		return tmpErr.WithError(err)
	default:
		return ErrUnknown.WithError(err)
	}
}

// GetHTTPErrorType 获取HTTP错误类型
func GetHTTPErrorType(err error) string {
	if err == nil {
		return ""
	}

	// 检查是否为上下文取消或超时
	if err == context.Canceled {
		return HTTPErrorTypeClient
	}
	if err == context.DeadlineExceeded {
		return HTTPErrorTypeTimeout
	}

	// 检查是否为HTTPError
	var httpErr *HTTPError
	if he, ok := err.(*HTTPError); ok {
		httpErr = he
	} else {
		// 非HTTPError，检查错误消息
		errMsg := err.Error()
		if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline") {
			return HTTPErrorTypeTimeout
		}
		if strings.Contains(errMsg, "connection") || strings.Contains(errMsg, "dial") {
			return HTTPErrorTypeConnection
		}
		return HTTPErrorTypeUnknown
	}

	// 根据状态码判断错误类型
	switch {
	case httpErr.StatusCode == 0:
		if strings.Contains(httpErr.Body, "timeout") || strings.Contains(httpErr.Body, "deadline") {
			return HTTPErrorTypeTimeout
		}
		return HTTPErrorTypeConnection
	case httpErr.StatusCode == http.StatusUnauthorized:
		return HTTPErrorTypeUnauthorized
	case httpErr.StatusCode == http.StatusForbidden:
		return HTTPErrorTypeForbidden
	case httpErr.StatusCode == http.StatusNotFound:
		return HTTPErrorTypeNotFound
	case httpErr.StatusCode >= 500:
		return HTTPErrorTypeServer
	case httpErr.StatusCode >= 400:
		return HTTPErrorTypeClient
	default:
		return HTTPErrorTypeUnknown
	}
}

// ParseJSONError 解析JSON错误响应
func ParseJSONError(data []byte) (int, string, error) {
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	err := json.Unmarshal(data, &resp)
	if err != nil {
		return 0, "", err
	}

	return resp.Code, resp.Message, nil
}
