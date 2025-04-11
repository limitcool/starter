// Package i18n 提供国际化支持
package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/log"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundle       *i18n.Bundle
	defaultLang  string
	localizer    map[string]*i18n.Localizer
	localizerMux sync.RWMutex
)

// Config 国际化配置
type Config struct {
	DefaultLanguage string   // 默认语言
	SupportLanguages []string // 支持的语言列表
}

// Setup 初始化i18n服务
func Setup(config Config, fsys fs.FS) error {
	// 设置默认语言
	if config.DefaultLanguage == "" {
		config.DefaultLanguage = "zh-CN"
	}
	defaultLang = config.DefaultLanguage

	// 创建bundle
	bundle = i18n.NewBundle(language.Make(defaultLang))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 创建localizer映射
	localizer = make(map[string]*i18n.Localizer)

	// 遍历文件系统中的所有翻译文件
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理json和toml文件
		if d.IsDir() || (!strings.HasSuffix(path, ".json") && !strings.HasSuffix(path, ".toml")) {
			return nil
		}

		ext := filepath.Ext(path)
		lang := strings.TrimSuffix(filepath.Base(path), ext)

		// 读取文件
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("无法读取语言文件 %s: %v", path, err)
		}

		// 加载到Bundle
		if _, err := bundle.ParseMessageFileBytes(data, path); err != nil {
			return fmt.Errorf("解析语言文件 %s 失败: %v", path, err)
		}

		log.Info("已加载语言文件", "path", path, "language", lang)

		// 初始化语言对应的localizer
		localizerMux.Lock()
		localizer[lang] = i18n.NewLocalizer(bundle, lang)
		localizerMux.Unlock()

		return nil
	})
}

// SetupWithEmbedFS 使用嵌入文件系统初始化i18n服务
func SetupWithEmbedFS(config Config, embedFS embed.FS, dir string) error {
	subFS, err := fs.Sub(embedFS, dir)
	if err != nil {
		return fmt.Errorf("无法访问嵌入文件系统目录 %s: %v", dir, err)
	}
	return Setup(config, subFS)
}

// Translate 翻译消息
func Translate(messageID string, lang string, args map[string]interface{}) string {
	if lang == "" {
		lang = defaultLang
	}

	localizerMux.RLock()
	loc, ok := localizer[lang]
	localizerMux.RUnlock()

	if !ok {
		// 如果指定的语言不存在，使用默认语言
		localizerMux.RLock()
		loc = localizer[defaultLang]
		localizerMux.RUnlock()
	}

	if loc == nil {
		// 如果没有可用的本地化器，返回消息ID
		return messageID
	}

	// 翻译消息
	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: args,
		DefaultMessage: &i18n.Message{
			ID:    messageID,
			Other: messageID, // 如果翻译不存在，返回消息ID本身
		},
	})

	if err != nil {
		log.Debug("翻译消息失败", "message_id", messageID, "lang", lang, "error", err)
		return messageID
	}

	return msg
}

// T 翻译消息的简便方法
func T(messageID string, lang string, args ...interface{}) string {
	// 将可变参数转换为map
	argsMap := make(map[string]interface{})

	if len(args) > 0 {
		// 如果是单个参数且为map，直接使用
		if len(args) == 1 {
			if m, ok := args[0].(map[string]interface{}); ok {
				argsMap = m
			}
		} else {
			// 偶数个参数作为键值对
			for i := 0; i < len(args); i += 2 {
				if i+1 < len(args) {
					key, ok := args[i].(string)
					if ok {
						argsMap[key] = args[i+1]
					}
				}
			}
		}
	}

	return Translate(messageID, lang, argsMap)
}

// GetSupportedLanguages 获取支持的语言列表
func GetSupportedLanguages() []string {
	localizerMux.RLock()
	defer localizerMux.RUnlock()

	languages := make([]string, 0, len(localizer))
	for lang := range localizer {
		languages = append(languages, lang)
	}

	return languages
}

// GetDefaultLanguage 获取默认语言
func GetDefaultLanguage() string {
	return defaultLang
}