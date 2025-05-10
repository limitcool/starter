package repository

import (
	"context"
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// MenuRepo 菜单仓库
type MenuRepo struct {
	DB          *gorm.DB
	GenericRepo Repository[model.Menu] // 使用接口而非具体实现
}

// GetMenuIDsByRoleID 获取角色关联的菜单ID
func (r *MenuRepo) GetMenuIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error) {
	// 创建RoleMenu的泛型仓库
	roleMenuRepo := NewGenericRepo[model.RoleMenu](r.DB)

	// 使用查询选项
	opts := &QueryOptions{
		Condition: "role_id = ?",
		Args:      []any{roleID},
	}

	// 获取所有关联记录
	roleMenus, err := roleMenuRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取角色菜单ID失败: roleID=%d", roleID))
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, roleMenu := range roleMenus {
		menuIDs = append(menuIDs, roleMenu.MenuID)
	}

	return menuIDs, nil
}

// GetMenusByIDs 根据ID列表获取菜单
func (r *MenuRepo) GetMenusByIDs(ctx context.Context, menuIDs []uint, preloads ...string) ([]*model.Menu, error) {
	if len(menuIDs) == 0 {
		return []*model.Menu{}, nil
	}

	// 使用泛型仓库的List方法
	opts := &QueryOptions{
		Condition: "id IN ? ORDER BY order_num",
		Args:      []any{menuIDs},
		Preloads:  preloads,
	}

	// 获取菜单（不分页，传入一个很大的页大小）
	menus, err := r.GenericRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, "查询菜单失败")
	}

	// 转换为指针切片
	menuPtrs := make([]*model.Menu, len(menus))
	for i := range menus {
		menuPtrs[i] = &menus[i]
	}

	return menuPtrs, nil
}

// NewMenuRepo 创建菜单仓库
func NewMenuRepo(db *gorm.DB) *MenuRepo {
	// 创建通用仓库并设置错误码
	genericRepo := NewGenericRepo[model.Menu](db).SetErrorCode(errorx.ErrNotFoundCodeValue)

	return &MenuRepo{
		DB:          db,
		GenericRepo: genericRepo,
	}
}

// GetByID 根据ID获取菜单
func (r *MenuRepo) GetByID(ctx context.Context, id uint) (*model.Menu, error) {
	// 使用仓库接口
	return r.GenericRepo.Get(ctx, id, nil)
}

// GetByIDWithRelations 根据ID获取菜单及其关联数据
// 使用预加载避免N+1查询问题
func (r *MenuRepo) GetByIDWithRelations(ctx context.Context, id uint) (*model.Menu, error) {
	// 使用泛型仓库的Get方法，并预加载关联
	opts := &QueryOptions{
		Preloads: []string{"Buttons", "APIs"},
	}

	// 查询菜单及其关联
	menu, err := r.GenericRepo.Get(ctx, id, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单及其关联失败: id=%d", id))
	}

	return menu, nil
}

// GetAll 获取所有菜单
// 使用预加载避免N+1查询问题
func (r *MenuRepo) GetAll(ctx context.Context) ([]*model.Menu, error) {
	// 使用泛型仓库的List方法，并预加载按钮
	opts := &QueryOptions{
		Preloads:  []string{"Buttons"},
		Condition: "1=1 ORDER BY order_num",
	}

	// 获取所有菜单（不分页，传入一个很大的页大小）
	menus, err := r.GenericRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, "查询所有菜单失败")
	}

	// 转换为指针切片
	menuPtrs := make([]*model.Menu, len(menus))
	for i := range menus {
		menuPtrs[i] = &menus[i]
	}

	return menuPtrs, nil
}

// Create 创建菜单
func (r *MenuRepo) Create(ctx context.Context, menu *model.Menu) error {
	// 使用仓库接口
	return r.GenericRepo.Create(ctx, menu)
}

// CreateButton 创建菜单按钮
func (r *MenuRepo) CreateButton(ctx context.Context, button *model.MenuButton) error {
	// 创建MenuButton的泛型仓库
	buttonRepo := NewGenericRepo[model.MenuButton](r.DB)
	return buttonRepo.Create(ctx, button)
}

// UpdateButton 更新菜单按钮
func (r *MenuRepo) UpdateButton(ctx context.Context, button *model.MenuButton) error {
	// 创建MenuButton的泛型仓库
	buttonRepo := NewGenericRepo[model.MenuButton](r.DB)
	return buttonRepo.Update(ctx, button)
}

// DeleteButton 删除菜单按钮
func (r *MenuRepo) DeleteButton(ctx context.Context, id uint) error {
	// 创建MenuButton的泛型仓库
	buttonRepo := NewGenericRepo[model.MenuButton](r.DB)
	return buttonRepo.Delete(ctx, id)
}

// GetButtonByID 根据ID获取菜单按钮
func (r *MenuRepo) GetButtonByID(ctx context.Context, id uint) (*model.MenuButton, error) {
	// 创建MenuButton的泛型仓库
	buttonRepo := NewGenericRepo[model.MenuButton](r.DB).SetErrorCode(errorx.ErrNotFoundCodeValue)
	return buttonRepo.Get(ctx, id, nil)
}

// Update 更新菜单
func (r *MenuRepo) Update(ctx context.Context, menu *model.Menu) error {
	// 使用仓库接口
	return r.GenericRepo.Update(ctx, menu)
}

// Delete 删除菜单
func (r *MenuRepo) Delete(ctx context.Context, id uint) error {
	// 检查是否有子菜单
	opts := &QueryOptions{
		Condition: "parent_id = ?",
		Args:      []any{id},
	}

	count, err := r.GenericRepo.Count(ctx, opts)
	if err != nil {
		return errorx.WrapError(err, "检查子菜单失败")
	}

	if count > 0 {
		return errorx.Errorf(errorx.ErrForbidden, "该菜单下有 %d 个子菜单，无法删除", count)
	}

	// 删除菜单
	return r.GenericRepo.Delete(ctx, id)
}

// GetByRoleID 获取角色菜单
// 使用预加载避免N+1查询问题
func (r *MenuRepo) GetByRoleID(ctx context.Context, roleID uint) ([]*model.Menu, error) {
	// 获取角色关联的菜单ID
	menuIDs, err := r.GetMenuIDsByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// 根据ID获取菜单
	return r.GetMenusByIDs(ctx, menuIDs, "Buttons")
}

// GetByUserID 获取用户菜单
// 使用JOIN查询和预加载避免N+1查询问题
func (r *MenuRepo) GetByUserID(ctx context.Context, userID uint) ([]*model.Menu, error) {
	// 创建UserRole的泛型仓库
	userRoleRepo := NewGenericRepo[model.UserRole](r.DB)

	// 使用查询选项
	userRoleOpts := &QueryOptions{
		Condition: "user_id = ?",
		Args:      []any{userID},
	}

	// 获取所有关联记录
	userRoles, err := userRoleRepo.List(ctx, 1, 1000, userRoleOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询用户角色关联失败: userID=%d", userID))
	}

	// 提取角色ID
	var roleIDs []uint
	for _, userRole := range userRoles {
		roleIDs = append(roleIDs, userRole.RoleID)
	}

	if len(roleIDs) == 0 {
		return []*model.Menu{}, nil
	}

	// 创建RoleMenu的泛型仓库
	roleMenuRepo := NewGenericRepo[model.RoleMenu](r.DB)

	// 使用查询选项
	roleMenuOpts := &QueryOptions{
		Condition: "role_id IN ?",
		Args:      []any{roleIDs},
	}

	// 获取所有关联记录
	roleMenus, err := roleMenuRepo.List(ctx, 1, 1000, roleMenuOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询角色菜单关联失败: roleIDs=%v", roleIDs))
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, roleMenu := range roleMenus {
		menuIDs = append(menuIDs, roleMenu.MenuID)
	}

	if len(menuIDs) == 0 {
		return []*model.Menu{}, nil
	}

	// 最后获取菜单，并过滤状态和类型
	opts := &QueryOptions{
		Condition: "id IN ? AND status = ? AND type IN ?",
		Args:      []any{menuIDs, 1, []int8{0, 1}},
		Preloads:  []string{"Buttons"},
	}

	// 获取菜单
	menus, err := r.GenericRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, "查询用户菜单失败")
	}

	// 转换为指针切片
	menuPtrs := make([]*model.Menu, len(menus))
	for i := range menus {
		menuPtrs[i] = &menus[i]
	}

	return menuPtrs, nil
}

// GetPermsByUserID 获取用户菜单权限
// 使用JOIN查询避免N+1查询问题
func (r *MenuRepo) GetPermsByUserID(ctx context.Context, userID uint) ([]string, error) {
	// 创建UserRole的泛型仓库
	userRoleRepo := NewGenericRepo[model.UserRole](r.DB)

	// 使用查询选项
	userRoleOpts := &QueryOptions{
		Condition: "user_id = ?",
		Args:      []any{userID},
	}

	// 获取所有关联记录
	userRoles, err := userRoleRepo.List(ctx, 1, 1000, userRoleOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询用户角色关联失败: userID=%d", userID))
	}

	// 提取角色ID
	var roleIDs []uint
	for _, userRole := range userRoles {
		roleIDs = append(roleIDs, userRole.RoleID)
	}

	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	// 使用角色ID获取权限
	return r.GetPermsByRoleIDs(ctx, roleIDs)
}

// GetPermsByRoleIDs 获取角色菜单权限
// 使用JOIN查询避免N+1查询问题
func (r *MenuRepo) GetPermsByRoleIDs(ctx context.Context, roleIDs []uint) ([]string, error) {
	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	// 创建RoleMenu的泛型仓库
	roleMenuRepo := NewGenericRepo[model.RoleMenu](r.DB)

	// 使用查询选项
	roleMenuOpts := &QueryOptions{
		Condition: "role_id IN ?",
		Args:      []any{roleIDs},
	}

	// 获取所有关联记录
	roleMenus, err := roleMenuRepo.List(ctx, 1, 1000, roleMenuOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询角色菜单关联失败: roleIDs=%v", roleIDs))
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, roleMenu := range roleMenus {
		menuIDs = append(menuIDs, roleMenu.MenuID)
	}

	if len(menuIDs) == 0 {
		return []string{}, nil
	}

	// 使用查询选项
	menuOpts := &QueryOptions{
		Condition: "id IN ? AND status = ? AND perms != ''",
		Args:      []any{menuIDs, 1},
	}

	// 获取所有菜单
	menus, err := r.GenericRepo.List(ctx, 1, 1000, menuOpts)
	if err != nil {
		return nil, errorx.WrapError(err, "查询菜单权限失败")
	}

	// 提取权限标识
	var perms []string
	for _, menu := range menus {
		if menu.Perms != "" {
			perms = append(perms, menu.Perms)
		}
	}

	return perms, nil
}

// AssociateAPI 关联菜单和API
func (r *MenuRepo) AssociateAPI(ctx context.Context, menuID uint, apiIDs []uint) error {
	// 使用GenericRepo的Transaction方法
	return r.GenericRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// 验证菜单是否存在
		menuRepo := r.GenericRepo.WithTx(tx)

		// 检查菜单是否存在
		_, err := menuRepo.Get(ctx, menuID, nil)
		if err != nil {
			return errorx.WrapError(err, fmt.Sprintf("查询菜单失败: id=%d", menuID))
		}

		// 创建MenuAPI的泛型仓库
		menuAPIRepo := NewGenericRepo[model.MenuAPI](tx)

		// 使用查询选项
		opts := &QueryOptions{
			Condition: "menu_id = ?",
			Args:      []any{menuID},
		}

		// 获取所有关联记录
		menuAPIs, err := menuAPIRepo.List(ctx, 1, 1000, opts)
		if err != nil {
			return errorx.WrapError(err, fmt.Sprintf("查询菜单API关联失败: menuID=%d", menuID))
		}

		// 删除所有关联记录
		for _, menuAPI := range menuAPIs {
			if err := menuAPIRepo.Delete(ctx, menuAPI.ID); err != nil {
				return errorx.WrapError(err, fmt.Sprintf("删除菜单API关联失败: menuID=%d, apiID=%d", menuID, menuAPI.APIID))
			}
		}

		// 添加新关联
		if len(apiIDs) > 0 {
			for _, apiID := range apiIDs {
				menuAPI := model.MenuAPI{
					MenuID: menuID,
					APIID:  apiID,
				}
				if err := menuAPIRepo.Create(ctx, &menuAPI); err != nil {
					return errorx.WrapError(err, fmt.Sprintf("创建菜单API关联失败: menuID=%d, apiID=%d", menuID, apiID))
				}
			}
		}

		return nil
	})
}

// AssignMenuToRole 为角色分配菜单
func (r *MenuRepo) AssignMenuToRole(ctx context.Context, roleID uint, menuIDs []uint) error {
	// 使用GenericRepo的Transaction方法
	return r.GenericRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// 创建RoleMenu的泛型仓库
		roleMenuRepo := NewGenericRepo[model.RoleMenu](tx)

		// 使用查询选项
		opts := &QueryOptions{
			Condition: "role_id = ?",
			Args:      []any{roleID},
		}

		// 获取所有关联记录
		roleMenus, err := roleMenuRepo.List(ctx, 1, 1000, opts)
		if err != nil {
			return errorx.WrapError(err, fmt.Sprintf("查询角色菜单关联失败: roleID=%d", roleID))
		}

		// 删除所有关联记录
		for _, roleMenu := range roleMenus {
			if err := roleMenuRepo.Delete(ctx, roleMenu.ID); err != nil {
				return errorx.WrapError(err, fmt.Sprintf("删除角色菜单关联失败: roleID=%d, menuID=%d", roleID, roleMenu.MenuID))
			}
		}

		// 添加新的角色菜单关联
		if len(menuIDs) > 0 {
			for _, menuID := range menuIDs {
				roleMenu := model.RoleMenu{
					RoleID: roleID,
					MenuID: menuID,
				}
				if err := roleMenuRepo.Create(ctx, &roleMenu); err != nil {
					return errorx.WrapError(err, fmt.Sprintf("创建角色菜单关联失败: roleID=%d, menuID=%d", roleID, menuID))
				}
			}
		}

		return nil
	})
}

// BuildMenuTree 构建菜单树
func (r *MenuRepo) BuildMenuTree(ctx context.Context, menus []*model.Menu) []*model.MenuTree {
	// 创建一个映射，用于快速查找菜单
	menuMap := make(map[uint]*model.MenuTree)

	// 将所有菜单转换为树节点
	for _, menu := range menus {
		menuMap[menu.ID] = &model.MenuTree{
			Menu:     menu,
			Children: []*model.MenuTree{},
		}
	}

	// 构建树结构
	var rootMenus []*model.MenuTree
	for _, menu := range menus {
		if menu.ParentID == 0 {
			// 根菜单
			rootMenus = append(rootMenus, menuMap[menu.ID])
		} else {
			// 子菜单
			if parent, ok := menuMap[menu.ParentID]; ok {
				parent.Children = append(parent.Children, menuMap[menu.ID])
			}
		}
	}

	return rootMenus
}

// AddMenuAPI 添加菜单API关联
func (r *MenuRepo) AddMenuAPI(ctx context.Context, menuID uint, apiID uint) error {
	// 创建MenuAPI的泛型仓库
	menuAPIRepo := NewGenericRepo[model.MenuAPI](r.DB)

	// 创建菜单API关联
	menuAPI := model.MenuAPI{
		MenuID: menuID,
		APIID:  apiID,
	}
	return menuAPIRepo.Create(ctx, &menuAPI)
}

// ClearMenuAPIs 清除菜单API关联
func (r *MenuRepo) ClearMenuAPIs(ctx context.Context, menuID uint) error {
	// 创建MenuAPI的泛型仓库
	menuAPIRepo := NewGenericRepo[model.MenuAPI](r.DB)

	// 使用查询选项
	opts := &QueryOptions{
		Condition: "menu_id = ?",
		Args:      []any{menuID},
	}

	// 获取所有关联记录
	menuAPIs, err := menuAPIRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("查询菜单API关联失败: menuID=%d", menuID))
	}

	// 删除所有关联记录
	for _, menuAPI := range menuAPIs {
		if err := menuAPIRepo.Delete(ctx, menuAPI.ID); err != nil {
			return errorx.WrapError(err, fmt.Sprintf("删除菜单API关联失败: menuID=%d, apiID=%d", menuID, menuAPI.APIID))
		}
	}

	return nil
}

// GetMenuIDsByAPIID 获取API关联的所有菜单ID
// 使用直接查询避免N+1查询问题
func (r *MenuRepo) GetMenuIDsByAPIID(ctx context.Context, apiID uint) ([]uint, error) {
	// 创建MenuAPI的泛型仓库
	menuAPIRepo := NewGenericRepo[model.MenuAPI](r.DB)

	// 使用查询选项
	opts := &QueryOptions{
		Condition: "api_id = ?",
		Args:      []any{apiID},
	}

	// 获取所有关联记录
	menuAPIs, err := menuAPIRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询API关联的菜单失败: apiID=%d", apiID))
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, menuAPI := range menuAPIs {
		menuIDs = append(menuIDs, menuAPI.MenuID)
	}

	return menuIDs, nil
}

// GetRolesByMenuID 获取拥有该菜单的所有角色
// 使用JOIN查询避免N+1查询问题
func (r *MenuRepo) GetRolesByMenuID(ctx context.Context, menuID uint) ([]*model.Role, error) {
	// 创建RoleMenu的泛型仓库
	roleMenuRepo := NewGenericRepo[model.RoleMenu](r.DB)

	// 使用查询选项
	roleMenuOpts := &QueryOptions{
		Condition: "menu_id = ?",
		Args:      []any{menuID},
	}

	// 获取所有关联记录
	roleMenus, err := roleMenuRepo.List(ctx, 1, 1000, roleMenuOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单角色关联失败: menuID=%d", menuID))
	}

	// 提取角色ID
	var roleIDs []uint
	for _, roleMenu := range roleMenus {
		roleIDs = append(roleIDs, roleMenu.RoleID)
	}

	if len(roleIDs) == 0 {
		return []*model.Role{}, nil
	}

	// 创建Role的泛型仓库
	roleRepo := NewGenericRepo[model.Role](r.DB)

	// 使用查询选项
	opts := &QueryOptions{
		Condition: "id IN ?",
		Args:      []any{roleIDs},
	}

	// 获取所有角色
	roles, err := roleRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单关联的角色失败: menuID=%d", menuID))
	}

	// 转换为指针切片
	rolePtrs := make([]*model.Role, len(roles))
	for i := range roles {
		rolePtrs[i] = &roles[i]
	}

	return rolePtrs, nil
}
