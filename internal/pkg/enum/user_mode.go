package enum

// UserMode 用户模式
type UserMode string

const (
	// UserModeSeparate 分离模式 - 使用user和admin_user两张表，完整的角色和权限管理
	UserModeSeparate UserMode = "separate"

	// UserModeSimple 简单模式 - 只使用user一张表，通过is_admin字段区分，不使用角色和权限
	UserModeSimple UserMode = "simple"
)

// IsValid 检查用户模式是否有效
func (m UserMode) IsValid() bool {
	return m == UserModeSeparate || m == UserModeSimple
}

// String 返回用户模式的字符串表示
func (m UserMode) String() string {
	return string(m)
}

// GetUserMode 根据字符串获取用户模式
func GetUserMode(mode string) UserMode {
	switch mode {
	case string(UserModeSeparate):
		return UserModeSeparate
	case string(UserModeSimple):
		return UserModeSimple
	case "unified", "simple_unified": // 兼容旧版本
		return UserModeSimple
	default:
		return UserModeSeparate // 默认使用分离模式
	}
}
