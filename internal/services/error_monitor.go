package services

import (
	"context"
	"sync"
	"time"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// ErrorStat 错误统计信息
type ErrorStat struct {
	Code      int       // 错误码
	Count     int       // 出现次数
	FirstSeen time.Time // 首次出现时间
	LastSeen  time.Time // 最后出现时间
}

// ErrorMonitorService 错误监控服务
type ErrorMonitorService struct {
	mu         sync.RWMutex
	errorStats map[int]*ErrorStat // 按错误码统计
	threshold  int                // 报警阈值
}

// NewErrorMonitorService 创建错误监控服务
func NewErrorMonitorService() *ErrorMonitorService {
	return &ErrorMonitorService{
		errorStats: make(map[int]*ErrorStat),
		threshold:  10, // 默认阈值：同一错误出现10次触发报警
	}
}

// RecordError 记录错误
func (s *ErrorMonitorService) RecordError(ctx context.Context, err error) {
	// 只处理AppError类型的错误
	if !errorx.IsAppErr(err) {
		return
	}

	appErr, ok := err.(*errorx.AppError)
	if !ok {
		return
	}
	code := appErr.GetErrorCode()
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	// 更新统计信息
	if stat, exists := s.errorStats[code]; exists {
		stat.Count++
		stat.LastSeen = now

		// 检查是否达到报警阈值
		if stat.Count == s.threshold {
			s.triggerAlert(ctx, appErr, stat)
		}
	} else {
		// 新错误
		s.errorStats[code] = &ErrorStat{
			Code:      code,
			Count:     1,
			FirstSeen: now,
			LastSeen:  now,
		}
	}
}

// GetErrorStats 获取错误统计信息
func (s *ErrorMonitorService) GetErrorStats(ctx context.Context) map[int]*ErrorStat {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 创建副本
	stats := make(map[int]*ErrorStat, len(s.errorStats))
	for k, v := range s.errorStats {
		stats[k] = &ErrorStat{
			Code:      v.Code,
			Count:     v.Count,
			FirstSeen: v.FirstSeen,
			LastSeen:  v.LastSeen,
		}
	}

	return stats
}

// SetThreshold 设置报警阈值
func (s *ErrorMonitorService) SetThreshold(ctx context.Context, threshold int) {
	if threshold <= 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.threshold = threshold
}

// ResetStats 重置统计信息
func (s *ErrorMonitorService) ResetStats(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.errorStats = make(map[int]*ErrorStat)
}

// triggerAlert 触发报警
func (s *ErrorMonitorService) triggerAlert(ctx context.Context, err *errorx.AppError, stat *ErrorStat) {
	// 这里可以实现各种报警方式，如发送邮件、短信、钉钉等
	// 目前仅记录日志
	logger.WarnContext(ctx, "Error threshold reached",
		"error_code", err.GetErrorCode(),
		"error_message", err.GetErrorMsg(),
		"error_count", stat.Count,
		"first_seen", stat.FirstSeen,
		"last_seen", stat.LastSeen,
	)
}
