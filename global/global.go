package global

type RequestIDType string
type ClientIpType string
type TokenType string

const (
	RequestIDKey RequestIDType = "request_id"
	ClientIp     ClientIpType  = "client_ip"
	Token        TokenType     = "token"
)

// 请使用 services.Instance().GetConfig() 获取应用配置
