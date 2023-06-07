package global

import "github.com/limitcool/starter/configs"

type RequestIDType string
type ClientIpType string
type TokenType string

const (
	RequestIDKey RequestIDType = "request_id"
	ClientIp     ClientIpType  = "client_ip"
	Token        TokenType     = "token"
)

var Config *configs.Config
