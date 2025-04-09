package enum

// MenuType 菜单类型
type MenuType int8

const (
	MenuTypeDirectory MenuType = iota // 目录
	MenuTypeMenu                      // 菜单
	MenuTypeButton                    // 按钮
)

// PermissionType 权限类型
type PermissionType int8

const (
	PermissionTypeMenu      PermissionType = iota + 1 // 菜单
	PermissionTypeOperation                           // 操作
	PermissionTypeAPI                                 // API
)

// UserType 用户类型
type UserType uint8

const (
	UserTypeSysUser UserType = iota + 1 // 系统用户
	UserTypeUser                        // 普通用户
)

func (u UserType) String() string {
	return []string{"sys_user", "user"}[u-1]
}

type TokenType uint8

const (
	TokenTypeAccess  TokenType = iota + 1 // 访问令牌
	TokenTypeRefresh                      // 刷新令牌
)

func (t TokenType) String() string {
	return []string{"access_token", "refresh_token"}[t-1]
}
