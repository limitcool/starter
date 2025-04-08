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
