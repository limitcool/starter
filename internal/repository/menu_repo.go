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

// MenuQueryBuilder 菜单查询构建器
type MenuQueryBuilder struct {
	db *gorm.DB
}

// NewMenuQueryBuilder 创建菜单查询构建器
func NewMenuQueryBuilder(db *gorm.DB) *MenuQueryBuilder {
	return &MenuQueryBuilder{db: db}
}

// WithPreload 添加预加载
func (qb *MenuQueryBuilder) WithPreload(preloads ...string) *MenuQueryBuilder {
	db := qb.db
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	qb.db = db
	return qb
}

// WithOrder 添加排序
func (qb *MenuQueryBuilder) WithOrder(order string) *MenuQueryBuilder {
	qb.db = qb.db.Order(order)
	return qb
}

// WithIDs 按ID查询
func (qb *MenuQueryBuilder) WithIDs(ids []uint) *MenuQueryBuilder {
	if len(ids) > 0 {
		qb.db = qb.db.Where("id IN ?", ids)
	}
	return qb
}

// WithStatus 按状态查询
func (qb *MenuQueryBuilder) WithStatus(status int) *MenuQueryBuilder {
	qb.db = qb.db.Where("status = ?", status)
	return qb
}

// WithTypes 按类型查询
func (qb *MenuQueryBuilder) WithTypes(types []int8) *MenuQueryBuilder {
	if len(types) > 0 {
		qb.db = qb.db.Where("type IN ?", types)
	}
	return qb
}

// Find 执行查询
func (qb *MenuQueryBuilder) Find(ctx context.Context, result any) error {
	return qb.db.WithContext(ctx).Find(result).Error
}

// First 查询第一条记录
func (qb *MenuQueryBuilder) First(ctx context.Context, result any) error {
	return qb.db.WithContext(ctx).First(result).Error
}

// Count 查询记录数
func (qb *MenuQueryBuilder) Count(ctx context.Context) (int64, error) {
	var count int64
	err := qb.db.WithContext(ctx).Count(&count).Error
	return count, err
}

// WithTransaction 使用事务执行函数
func (r *MenuRepo) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	// 开始事务
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 执行函数
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return errorx.WrapError(err, "提交事务失败")
	}

	return nil
}

// GetMenuIDsByRoleID 获取角色关联的菜单ID
func (r *MenuRepo) GetMenuIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error) {
	var menuIDs []uint
	err := r.DB.WithContext(ctx).
		Model(&model.RoleMenu{}).
		Select("menu_id").
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDs).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取角色菜单ID失败: roleID=%d", roleID))
	}

	return menuIDs, nil
}

// GetMenusByIDs 根据ID列表获取菜单
func (r *MenuRepo) GetMenusByIDs(ctx context.Context, menuIDs []uint, preloads ...string) ([]*model.Menu, error) {
	if len(menuIDs) == 0 {
		return []*model.Menu{}, nil
	}

	var menus []*model.Menu
	qb := NewMenuQueryBuilder(r.DB.Model(&model.Menu{}))

	// 添加预加载
	if len(preloads) > 0 {
		qb = qb.WithPreload(preloads...)
	}

	// 执行查询
	err := qb.WithIDs(menuIDs).
		WithOrder("order_num").
		Find(ctx, &menus)

	if err != nil {
		return nil, errorx.WrapError(err, "查询菜单失败")
	}

	return menus, nil
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
// 使用预加载避免N+1查询问题
func (r *MenuRepo) GetByIDWithRelations(ctx context.Context, id uint) (*model.Menu, error) {
	var menu model.Menu

	// 使用预加载一次性加载所有关联数据
	err := r.DB.WithContext(ctx).
		Preload("Buttons"). // 预加载按钮
		Preload("APIs").    // 预加载API
		Where("id = ?", id).
		First(&menu).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		notFoundErr := errorx.Errorf(errorx.ErrNotFound, "菜单ID %d 不存在", id)
		return nil, errorx.WrapError(notFoundErr, "")
	}

	if err != nil {
		return nil, errorx.WrapError(err, "查询菜单失败")
	}

	return &menu, nil
}

// GetAll 获取所有菜单
// 使用预加载避免N+1查询问题
func (r *MenuRepo) GetAll(ctx context.Context) ([]*model.Menu, error) {
	var menus []*model.Menu

	// 使用预加载一次性加载所有菜单及其按钮
	err := r.DB.WithContext(ctx).
		Preload("Buttons"). // 预加载按钮
		Order("order_num").
		Find(&menus).Error

	if err != nil {
		return nil, errorx.WrapError(err, "查询所有菜单失败")
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
	// 使用JOIN查询直接获取菜单ID
	var menuIDs []uint
	err := r.DB.WithContext(ctx).
		Table("sys_menu").
		Select("DISTINCT sys_menu.id").
		Joins("JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_menu.role_id").
		Where("sys_user_role.user_id = ? AND sys_menu.status = ? AND sys_menu.type IN ?", userID, 1, []int8{0, 1}).
		Pluck("sys_menu.id", &menuIDs).Error

	if err != nil {
		return nil, errorx.WrapError(err, "查询用户菜单ID失败")
	}

	// 根据ID获取菜单
	return r.GetMenusByIDs(ctx, menuIDs, "Buttons")
}

// GetPermsByUserID 获取用户菜单权限
// 使用JOIN查询避免N+1查询问题
func (r *MenuRepo) GetPermsByUserID(ctx context.Context, userID uint) ([]string, error) {
	// 使用单个JOIN查询直接获取权限标识
	var perms []string
	err := r.DB.WithContext(ctx).
		Model(&model.Menu{}).
		Select("DISTINCT sys_menu.perms").
		Joins("JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_menu.role_id").
		Where("sys_user_role.user_id = ? AND sys_menu.status = ? AND sys_menu.perms != ''", userID, 1).
		Pluck("perms", &perms).Error

	if err != nil {
		return nil, errorx.WrapError(err, "查询用户菜单权限失败")
	}

	return perms, nil
}

// GetPermsByRoleIDs 获取角色菜单权限
// 使用JOIN查询避免N+1查询问题
func (r *MenuRepo) GetPermsByRoleIDs(ctx context.Context, roleIDs []uint) ([]string, error) {
	// 使用JOIN查询直接获取权限标识
	var perms []string
	err := r.DB.WithContext(ctx).
		Model(&model.Menu{}).
		Select("DISTINCT sys_menu.perms").
		Joins("JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Where("sys_role_menu.role_id IN ? AND sys_menu.status = ? AND sys_menu.perms != ''", roleIDs, 1).
		Pluck("perms", &perms).Error

	if err != nil {
		return nil, errorx.WrapError(err, "查询角色菜单权限失败")
	}

	return perms, nil
}

// AssociateAPI 关联菜单和API
func (r *MenuRepo) AssociateAPI(ctx context.Context, menuID uint, apiIDs []uint) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		// 获取菜单
		var menu model.Menu
		if err := tx.WithContext(ctx).First(&menu, menuID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorx.Errorf(errorx.ErrNotFound, "菜单ID %d 不存在", menuID)
			}
			return errorx.WrapError(err, "查询菜单失败")
		}

		// 删除原有关联
		if err := tx.WithContext(ctx).Where("menu_id = ?", menuID).Delete(&model.MenuAPI{}).Error; err != nil {
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
			if err := tx.WithContext(ctx).Create(&menuAPIs).Error; err != nil {
				return errorx.WrapError(err, "创建菜单API关联失败")
			}
		}

		return nil
	})
}

// AssignMenuToRole 为角色分配菜单
func (r *MenuRepo) AssignMenuToRole(ctx context.Context, roleID uint, menuIDs []uint) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		// 删除原有的角色菜单关联
		if err := tx.WithContext(ctx).Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
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
				return errorx.WrapError(err, "创建角色菜单关联失败")
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

// CreateAssociation 创建关联记录
func (r *MenuRepo) CreateAssociation(ctx context.Context, model any) error {
	err := r.DB.WithContext(ctx).Create(model).Error
	if err != nil {
		return errorx.WrapError(err, "创建关联记录失败")
	}
	return nil
}

// DeleteAssociation 删除关联记录
func (r *MenuRepo) DeleteAssociation(ctx context.Context, model any, query any, args ...any) error {
	err := r.DB.WithContext(ctx).Where(query, args...).Delete(model).Error
	if err != nil {
		return errorx.WrapError(err, "删除关联记录失败")
	}
	return nil
}

// AddMenuAPI 添加菜单API关联
func (r *MenuRepo) AddMenuAPI(ctx context.Context, menuID uint, apiID uint) error {
	menuAPI := model.MenuAPI{
		MenuID: menuID,
		APIID:  apiID,
	}
	return r.CreateAssociation(ctx, &menuAPI)
}

// ClearMenuAPIs 清除菜单API关联
func (r *MenuRepo) ClearMenuAPIs(ctx context.Context, menuID uint) error {
	return r.DeleteAssociation(ctx, &model.MenuAPI{}, "menu_id = ?", menuID)
}

// GetMenuIDsByAPIID 获取API关联的所有菜单ID
// 使用直接查询避免N+1查询问题
func (r *MenuRepo) GetMenuIDsByAPIID(ctx context.Context, apiID uint) ([]uint, error) {
	// 直接查询菜单ID，避免中间对象转换
	var menuIDs []uint
	err := r.DB.WithContext(ctx).
		Model(&model.MenuAPI{}).
		Select("menu_id").
		Where("api_id = ?", apiID).
		Pluck("menu_id", &menuIDs).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询API关联的菜单失败: apiID=%d", apiID))
	}

	return menuIDs, nil
}

// GetRolesByMenuID 获取拥有该菜单的所有角色
// 使用JOIN查询避免N+1查询问题
func (r *MenuRepo) GetRolesByMenuID(ctx context.Context, menuID uint) ([]*model.Role, error) {
	// 使用JOIN查询直接获取角色信息
	var roles []*model.Role
	err := r.DB.WithContext(ctx).
		Joins("JOIN sys_role_menu ON sys_role_menu.role_id = sys_role.id").
		Where("sys_role_menu.menu_id = ?", menuID).
		Find(&roles).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单关联的角色失败: menuID=%d", menuID))
	}

	return roles, nil
}
