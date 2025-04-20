package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// MenuRepo 菜单仓库
type MenuRepo struct {
	DB          *gorm.DB
	genericRepo *GenericRepo[model.Menu] // 泛型仓库
}

// NewMenuRepo 创建菜单仓库
func NewMenuRepo(db *gorm.DB) *MenuRepo {
	genericRepo := NewGenericRepo[model.Menu](db)
	genericRepo.SetErrorCode(errorx.ErrorNotFoundCode) // 设置错误码

	return &MenuRepo{
		DB:          db,
		genericRepo: genericRepo,
	}
}

// GetByID 根据ID获取菜单
func (r *MenuRepo) GetByID(ctx context.Context, id uint) (*model.Menu, error) {
	// 使用泛型仓库
	return r.genericRepo.GetByID(ctx, id)
}

// GetByIDWithRelations 根据ID获取菜单及其关联数据
func (r *MenuRepo) GetByIDWithRelations(ctx context.Context, id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&menu).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		notFoundErr := errorx.Errorf(errorx.ErrNotFound, "菜单ID %d 不存在", id)
		return nil, errorx.WrapError(notFoundErr, "")
	}

	if err != nil {
		return nil, errorx.WrapError(err, "查询菜单失败")
	}

	// 加载按钮
	if err := r.DB.WithContext(ctx).Where("menu_id = ?", menu.ID).Find(&menu.Buttons).Error; err != nil {
		return nil, errorx.WrapError(err, "查询菜单按钮失败")
	}

	// 加载API
	if err := r.DB.WithContext(ctx).Model(&menu).Association("APIs").Find(&menu.APIs); err != nil {
		return nil, errorx.WrapError(err, "查询菜单关联API失败")
	}

	return &menu, nil
}

// GetAll 获取所有菜单
func (r *MenuRepo) GetAll(ctx context.Context) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.DB.WithContext(ctx).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询所有菜单失败")
	}

	// 加载按钮
	for _, menu := range menus {
		if err := r.DB.WithContext(ctx).Where("menu_id = ?", menu.ID).Find(&menu.Buttons).Error; err != nil {
			return nil, errorx.WrapError(err, "查询菜单按钮失败")
		}
	}

	return menus, nil
}

// Create 创建菜单
func (r *MenuRepo) Create(ctx context.Context, menu *model.Menu) error {
	// 使用泛型仓库
	return r.genericRepo.Create(ctx, menu)
}

// CreateButton 创建菜单按钮
func (r *MenuRepo) CreateButton(ctx context.Context, button *model.MenuButton) error {
	err := r.DB.WithContext(ctx).Create(button).Error
	if err != nil {
		return errorx.WrapError(err, "创建菜单按钮失败")
	}
	return nil
}

// UpdateButton 更新菜单按钮
func (r *MenuRepo) UpdateButton(ctx context.Context, button *model.MenuButton) error {
	err := r.DB.WithContext(ctx).Save(button).Error
	if err != nil {
		return errorx.WrapError(err, "更新菜单按钮失败")
	}
	return nil
}

// DeleteButton 删除菜单按钮
func (r *MenuRepo) DeleteButton(ctx context.Context, id uint) error {
	err := r.DB.WithContext(ctx).Delete(&model.MenuButton{}, id).Error
	if err != nil {
		return errorx.WrapError(err, "删除菜单按钮失败")
	}
	return nil
}

// GetButtonByID 根据ID获取菜单按钮
func (r *MenuRepo) GetButtonByID(ctx context.Context, id uint) (*model.MenuButton, error) {
	var button model.MenuButton
	err := r.DB.WithContext(ctx).First(&button, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		notFoundErr := errorx.Errorf(errorx.ErrNotFound, "菜单按钮ID %d 不存在", id)
		return nil, errorx.WrapError(notFoundErr, "")
	}

	if err != nil {
		return nil, errorx.WrapError(err, "查询菜单按钮失败")
	}

	return &button, nil
}

// Update 更新菜单
func (r *MenuRepo) Update(ctx context.Context, menu *model.Menu) error {
	// 使用泛型仓库
	return r.genericRepo.Update(ctx, menu)
}

// Delete 删除菜单
func (r *MenuRepo) Delete(ctx context.Context, id uint) error {
	// 检查是否有子菜单
	var count int64
	if err := r.DB.WithContext(ctx).Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return errorx.WrapError(err, "检查子菜单失败")
	}

	if count > 0 {
		return errorx.Errorf(errorx.ErrForbidden, "该菜单下有 %d 个子菜单，无法删除", count)
	}

	// 删除菜单
	if err := r.DB.WithContext(ctx).Delete(&model.Menu{}, id).Error; err != nil {
		return errorx.WrapError(err, "删除菜单失败")
	}

	return nil
}

// GetByRoleID 获取角色菜单
func (r *MenuRepo) GetByRoleID(ctx context.Context, roleID uint) ([]*model.Menu, error) {
	// 查询角色关联的菜单ID
	var roleMenus []model.RoleMenu
	err := r.DB.WithContext(ctx).Where("role_id = ?", roleID).Find(&roleMenus).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询角色菜单关联失败")
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	if len(menuIDs) == 0 {
		return []*model.Menu{}, nil
	}

	// 查询菜单
	var menus []*model.Menu
	err = r.DB.WithContext(ctx).Where("id IN ?", menuIDs).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询菜单失败")
	}

	return menus, nil
}

// GetByUserID 获取用户菜单
func (r *MenuRepo) GetByUserID(ctx context.Context, userID uint) ([]*model.Menu, error) {
	// 1. 获取用户角色
	var userRoles []model.UserRole
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询用户角色失败")
	}

	// 提取角色ID
	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	if len(roleIDs) == 0 {
		return []*model.Menu{}, nil
	}

	// 2. 获取角色关联的菜单ID
	var roleMenus []model.RoleMenu
	err = r.DB.WithContext(ctx).Where("role_id IN ?", roleIDs).Find(&roleMenus).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询角色菜单关联失败")
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	if len(menuIDs) == 0 {
		return []*model.Menu{}, nil
	}

	// 3. 查询菜单信息
	var menus []*model.Menu
	err = r.DB.WithContext(ctx).Where("id IN ? AND status = ? AND type IN ?", menuIDs, 1, []int8{0, 1}).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询菜单失败")
	}

	return menus, nil
}

// GetPermsByUserID 获取用户菜单权限
func (r *MenuRepo) GetPermsByUserID(ctx context.Context, userID uint) ([]string, error) {
	// 获取用户角色
	var userRoles []model.UserRole
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询用户角色失败")
	}

	// 提取角色ID
	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	return r.GetPermsByRoleIDs(ctx, roleIDs)
}

// GetPermsByRoleIDs 获取角色菜单权限
func (r *MenuRepo) GetPermsByRoleIDs(ctx context.Context, roleIDs []uint) ([]string, error) {
	// 查询角色菜单
	var roleMenus []model.RoleMenu
	err := r.DB.WithContext(ctx).Where("role_id IN ?", roleIDs).Find(&roleMenus).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询角色菜单关联失败")
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	if len(menuIDs) == 0 {
		return []string{}, nil
	}

	// 查询菜单权限标识
	var perms []string
	err = r.DB.WithContext(ctx).Model(&model.Menu{}).
		Where("id IN ? AND status = ? AND perms != ''", menuIDs, 1).
		Pluck("perms", &perms).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询菜单权限失败")
	}

	return perms, nil
}

// AssociateAPI 关联菜单和API
func (r *MenuRepo) AssociateAPI(menuID uint, apiIDs []uint) error {
	// 开始事务
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 获取菜单
	var menu model.Menu
	if err := tx.First(&menu, menuID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.Errorf(errorx.ErrNotFound, "菜单ID %d 不存在", menuID)
		}
		return errorx.WrapError(err, "查询菜单失败")
	}

	// 删除原有关联
	if err := tx.Where("menu_id = ?", menuID).Delete(&model.MenuAPI{}).Error; err != nil {
		tx.Rollback()
		return errorx.WrapError(err, "删除菜单API关联失败")
	}

	// 添加新关联
	if len(apiIDs) > 0 {
		var menuAPIs []model.MenuAPI
		for _, apiID := range apiIDs {
			menuAPIs = append(menuAPIs, model.MenuAPI{
				MenuID: menuID,
				APIID:  apiID,
			})
		}
		if err := tx.Create(&menuAPIs).Error; err != nil {
			tx.Rollback()
			return errorx.WrapError(err, "创建菜单API关联失败")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errorx.WrapError(err, "提交事务失败")
	}

	return nil
}

// AssignMenuToRole 为角色分配菜单
func (r *MenuRepo) AssignMenuToRole(ctx context.Context, roleID uint, menuIDs []uint) error {
	// 开始事务
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的角色菜单关联
	if err := tx.WithContext(ctx).Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return errorx.WrapError(err, "删除角色菜单关联失败")
	}

	// 添加新的角色菜单关联
	if len(menuIDs) > 0 {
		var roleMenus []model.RoleMenu
		for _, menuID := range menuIDs {
			roleMenus = append(roleMenus, model.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			})
		}
		if err := tx.WithContext(ctx).Create(&roleMenus).Error; err != nil {
			tx.Rollback()
			return errorx.WrapError(err, "创建角色菜单关联失败")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errorx.WrapError(err, "提交事务失败")
	}

	return nil
}

// BuildMenuTree 构建菜单树
func (r *MenuRepo) BuildMenuTree(menus []*model.Menu) []*model.MenuTree {
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
	menuAPI := model.MenuAPI{
		MenuID: menuID,
		APIID:  apiID,
	}
	err := r.DB.WithContext(ctx).Create(&menuAPI).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("添加菜单API关联失败: menuID=%d, apiID=%d", menuID, apiID))
	}
	return nil
}

// ClearMenuAPIs 清除菜单API关联
func (r *MenuRepo) ClearMenuAPIs(ctx context.Context, menuID uint) error {
	err := r.DB.WithContext(ctx).Where("menu_id = ?", menuID).Delete(&model.MenuAPI{}).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("清除菜单API关联失败: menuID=%d", menuID))
	}
	return nil
}

// GetMenuIDsByAPIID 获取API关联的所有菜单ID
func (r *MenuRepo) GetMenuIDsByAPIID(ctx context.Context, apiID uint) ([]uint, error) {
	var menuAPIs []model.MenuAPI
	err := r.DB.WithContext(ctx).Where("api_id = ?", apiID).Find(&menuAPIs).Error
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询API关联的菜单失败: apiID=%d", apiID))
	}

	var menuIDs []uint
	for _, ma := range menuAPIs {
		menuIDs = append(menuIDs, ma.MenuID)
	}

	return menuIDs, nil
}

// GetRolesByMenuID 获取拥有该菜单的所有角色
func (r *MenuRepo) GetRolesByMenuID(ctx context.Context, menuID uint) ([]*model.Role, error) {
	var roleMenus []model.RoleMenu
	err := r.DB.WithContext(ctx).Where("menu_id = ?", menuID).Find(&roleMenus).Error
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单角色关联失败: menuID=%d", menuID))
	}

	var roleIDs []uint
	for _, rm := range roleMenus {
		roleIDs = append(roleIDs, rm.RoleID)
	}

	if len(roleIDs) == 0 {
		return []*model.Role{}, nil
	}

	var roles []*model.Role
	err = r.DB.Where("id IN ?", roleIDs).Find(&roles).Error
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询角色失败: roleIDs=%v", roleIDs))
	}
	return roles, nil
}
