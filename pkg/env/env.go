package env

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

type Environment string

const (
	Dev  Environment = "dev"
	Test Environment = "test"
	Prod Environment = "prod"
)

// Get 获取当前环境
func Get() Environment {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return Dev // 默认为开发环境
	}
	switch strings.ToLower(env) {
	case "dev", "development":
		return Dev
	case "test", "testing":
		return Test
	case "prod", "production":
		return Prod
	default:
		log.Warnf("未知环境: %s, 使用默认开发环境", env)
		return Dev
	}
}

// IsDev 是否为开发环境
func IsDev() bool {
	return Get() == Dev
}

// IsTest 是否为测试环境
func IsTest() bool {
	return Get() == Test
}

// IsProd 是否为生产环境
func IsProd() bool {
	return Get() == Prod
}

// String 返回环境名称
func (e Environment) String() string {
	return string(e)
}