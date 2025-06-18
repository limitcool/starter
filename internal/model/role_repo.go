package model

import (
	"context"

	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// RoleRepo 角色仓库
type RoleRepo struct {
	Repository[Role]
	db *gorm.DB
}

// NewRoleRepo 创建角色仓库
func NewRoleRepo(db *gorm.DB) *RoleRepo {
	genericRepo := NewGenericRepo[Role](db)
	genericRepo.ErrorCode = errorx.ErrorNotFoundCode

	return &RoleRepo{
		Repository: genericRepo,
		db:         db,
	}
}

// GetByName 根据名称获取角色
func (r *RoleRepo) GetByName(ctx context.Context, name string) (*Role, error) {
	opts := &QueryOptions{
		Condition: "name = ?",
		Args:      []any{name},
	}
	return r.Get(ctx, nil, opts)
}

// GetByKey 根据角色Key获取角色
func (r *RoleRepo) GetByKey(ctx context.Context, key string) (*Role, error) {
	opts := &QueryOptions{
		Condition: "key = ? AND status = ?",
		Args:      []any{key, 1},
	}
	return r.Get(ctx, nil, opts)
}

// GetEnabledRoles 获取启用的角色列表
func (r *RoleRepo) GetEnabledRoles(ctx context.Context) ([]Role, error) {
	opts := &QueryOptions{
		Condition: "status = ?",
		Args:      []any{1},
	}
	return r.List(ctx, 1, 1000, opts)
}

// ListWithPagination 获取角色列表（带分页和总数）
func (r *RoleRepo) ListWithPagination(ctx context.Context, page, pageSize int, opts *QueryOptions) ([]Role, int64, error) {
	var roles []Role
	var total int64

	// 获取总数
	countQuery := r.db.WithContext(ctx).Model(&Role{})
	if opts != nil && opts.Condition != "" {
		countQuery = countQuery.Where(opts.Condition, opts.Args...)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, errorx.WrapError(err, "获取角色总数失败")
	}

	// 获取数据
	query := r.db.WithContext(ctx).Model(&Role{})
	if opts != nil {
		if opts.Condition != "" {
			query = query.Where(opts.Condition, opts.Args...)
		}
	}

	// 应用分页
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, errorx.WrapError(err, "获取角色列表失败")
	}

	return roles, total, nil
}

// GetRolesByUserID 获取用户的角色列表
func (r *RoleRepo) GetRolesByUserID(ctx context.Context, userID int64) ([]Role, error) {
	var roles []Role
	err := r.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.status = ?", userID, 1).
		Find(&roles).Error

	if err != nil {
		return nil, errorx.WrapError(err, "获取用户角色失败")
	}

	return roles, nil
}

// 注意：在基于Casbin的设计中，角色权限分配由Casbin管理，不需要数据库关联表

// 注意：在基于Casbin的设计中，不需要角色直接关联菜单

// GetPermissionsByRoleID 获取角色的权限列表
func (r *RoleRepo) GetPermissionsByRoleID(ctx context.Context, roleID uint) ([]Permission, error) {
	var permissions []Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND permissions.status = ?", roleID, 1).
		Find(&permissions).Error

	if err != nil {
		return nil, errorx.WrapError(err, "获取角色权限失败")
	}

	return permissions, nil
}

// GetMenusByRoleID 获取角色的菜单列表
func (r *RoleRepo) GetMenusByRoleID(ctx context.Context, roleID uint) ([]Menu, error) {
	var menus []Menu
	err := r.db.WithContext(ctx).
		Table("menus").
		Joins("JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id = ? AND menus.status = ?", roleID, 1).
		Order("menus.sort ASC, menus.id ASC").
		Find(&menus).Error

	if err != nil {
		return nil, errorx.WrapError(err, "获取角色菜单失败")
	}

	return menus, nil
}

// PermissionRepo 权限仓库
type PermissionRepo struct {
	Repository[Permission]
	db *gorm.DB
}

// NewPermissionRepo 创建权限仓库
func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	genericRepo := NewGenericRepo[Permission](db)
	genericRepo.ErrorCode = errorx.ErrorNotFoundCode

	return &PermissionRepo{
		Repository: genericRepo,
		db:         db,
	}
}

// GetByName 根据名称获取权限
func (r *PermissionRepo) GetByName(ctx context.Context, name string) (*Permission, error) {
	opts := &QueryOptions{
		Condition: "name = ?",
		Args:      []any{name},
	}
	return r.Get(ctx, nil, opts)
}

// GetByKey 根据权限Key获取权限
func (r *PermissionRepo) GetByKey(ctx context.Context, key string) (*Permission, error) {
	opts := &QueryOptions{
		Condition: "key = ?",
		Args:      []any{key},
	}
	return r.Get(ctx, nil, opts)
}

// HasChildren 检查权限是否有子权限
func (r *PermissionRepo) HasChildren(ctx context.Context, parentID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Permission{}).Where("parent_id = ?", parentID).Count(&count).Error
	if err != nil {
		return false, errorx.WrapError(err, "检查子权限失败")
	}
	return count > 0, nil
}

// GetAll 获取所有权限
func (r *PermissionRepo) GetAll(ctx context.Context) ([]Permission, error) {
	var permissions []Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	if err != nil {
		return nil, errorx.WrapError(err, "获取所有权限失败")
	}
	return permissions, nil
}

// ListWithPagination 获取权限列表（带分页和总数）
func (r *PermissionRepo) ListWithPagination(ctx context.Context, page, pageSize int, opts *QueryOptions) ([]Permission, int64, error) {
	var permissions []Permission
	var total int64

	// 获取总数
	countQuery := r.db.WithContext(ctx).Model(&Permission{})
	if opts != nil && opts.Condition != "" {
		countQuery = countQuery.Where(opts.Condition, opts.Args...)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, errorx.WrapError(err, "获取权限总数失败")
	}

	// 获取数据
	query := r.db.WithContext(ctx).Model(&Permission{})
	if opts != nil {
		if opts.Condition != "" {
			query = query.Where(opts.Condition, opts.Args...)
		}
	}

	// 应用分页
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&permissions).Error; err != nil {
		return nil, 0, errorx.WrapError(err, "获取权限列表失败")
	}

	return permissions, total, nil
}

// GetByResource 根据资源和操作获取权限
func (r *PermissionRepo) GetByResource(ctx context.Context, resource, action string) (*Permission, error) {
	opts := &QueryOptions{
		Condition: "resource = ? AND action = ?",
		Args:      []any{resource, action},
	}
	return r.Get(ctx, nil, opts)
}

// GetEnabledPermissions 获取启用的权限列表
func (r *PermissionRepo) GetEnabledPermissions(ctx context.Context) ([]Permission, error) {
	opts := &QueryOptions{
		Condition: "status = ?",
		Args:      []any{1},
	}
	return r.List(ctx, 1, 1000, opts)
}

// GetPermissionsByUserID 获取用户的权限列表（通过角色）
func (r *PermissionRepo) GetPermissionsByUserID(ctx context.Context, userID int64) ([]Permission, error) {
	var permissions []Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.status = ?", userID, 1).
		Group("permissions.id").
		Find(&permissions).Error

	if err != nil {
		return nil, errorx.WrapError(err, "获取用户权限失败")
	}

	return permissions, nil
}

// MenuRepo 菜单仓库
type MenuRepo struct {
	Repository[Menu]
	db *gorm.DB
}

// NewMenuRepo 创建菜单仓库
func NewMenuRepo(db *gorm.DB) *MenuRepo {
	genericRepo := NewGenericRepo[Menu](db)
	genericRepo.ErrorCode = errorx.ErrorNotFoundCode

	return &MenuRepo{
		Repository: genericRepo,
		db:         db,
	}
}

// GetByName 根据名称获取菜单
func (r *MenuRepo) GetByName(ctx context.Context, name string) (*Menu, error) {
	opts := &QueryOptions{
		Condition: "name = ?",
		Args:      []any{name},
	}
	return r.Get(ctx, nil, opts)
}

// GetEnabledMenus 获取启用的菜单列表
func (r *MenuRepo) GetEnabledMenus(ctx context.Context) ([]Menu, error) {
	opts := &QueryOptions{
		Condition: "status = ?",
		Args:      []any{1},
	}
	return r.List(ctx, 1, 1000, opts)
}

// GetMenusByUserID 获取用户的菜单列表（通过角色）
func (r *MenuRepo) GetMenusByUserID(ctx context.Context, userID int64) ([]Menu, error) {
	var menus []Menu
	err := r.db.WithContext(ctx).
		Table("menus").
		Joins("JOIN role_menus ON menus.id = role_menus.menu_id").
		Joins("JOIN user_roles ON role_menus.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND menus.status = ?", userID, 1).
		Group("menus.id").
		Order("menus.sort ASC, menus.id ASC").
		Find(&menus).Error

	if err != nil {
		return nil, errorx.WrapError(err, "获取用户菜单失败")
	}

	return menus, nil
}

// GetUserMenuTree 获取用户菜单树
func (r *MenuRepo) GetUserMenuTree(ctx context.Context, userID int64) ([]Menu, error) {
	menus, err := r.GetMenusByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return r.BuildMenuTree(menus, 0), nil
}

// BuildMenuTree 构建菜单树
func (r *MenuRepo) BuildMenuTree(menus []Menu, parentID uint) []Menu {
	var tree []Menu

	for i, menu := range menus {
		if menu.ParentID == parentID {
			children := r.BuildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				menus[i].Children = children
			}
			tree = append(tree, menus[i])
		}
	}

	return tree
}

// HasChildren 检查菜单是否有子菜单
func (r *MenuRepo) HasChildren(ctx context.Context, parentID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Menu{}).Where("parent_id = ?", parentID).Count(&count).Error
	if err != nil {
		return false, errorx.WrapError(err, "检查子菜单失败")
	}
	return count > 0, nil
}

// GetAll 获取所有菜单
func (r *MenuRepo) GetAll(ctx context.Context, opts *QueryOptions) ([]Menu, error) {
	var menus []Menu
	query := r.db.WithContext(ctx).Model(&Menu{})

	if opts != nil {
		if opts.Condition != "" {
			query = query.Where(opts.Condition, opts.Args...)
		}
	}

	// 默认按排序字段排序
	query = query.Order("sort_order ASC, id ASC")

	err := query.Find(&menus).Error
	if err != nil {
		return nil, errorx.WrapError(err, "获取所有菜单失败")
	}
	return menus, nil
}

// ListWithPagination 获取菜单列表（带分页和总数）
func (r *MenuRepo) ListWithPagination(ctx context.Context, page, pageSize int, opts *QueryOptions) ([]Menu, int64, error) {
	var menus []Menu
	var total int64

	// 获取总数
	countQuery := r.db.WithContext(ctx).Model(&Menu{})
	if opts != nil && opts.Condition != "" {
		countQuery = countQuery.Where(opts.Condition, opts.Args...)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, errorx.WrapError(err, "获取菜单总数失败")
	}

	// 获取数据
	query := r.db.WithContext(ctx).Model(&Menu{})
	if opts != nil {
		if opts.Condition != "" {
			query = query.Where(opts.Condition, opts.Args...)
		}
	}

	// 应用分页和排序
	query = query.Order("sort_order ASC, id ASC")
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&menus).Error; err != nil {
		return nil, 0, errorx.WrapError(err, "获取菜单列表失败")
	}

	return menus, total, nil
}

// UpdateSort 批量更新菜单排序
func (r *MenuRepo) UpdateSort(ctx context.Context, menuSorts []dto.MenuSort) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, sort := range menuSorts {
			if err := tx.Model(&Menu{}).Where("id = ?", sort.ID).Update("sort_order", sort.SortOrder).Error; err != nil {
				return errorx.WrapError(err, "更新菜单排序失败")
			}
		}
		return nil
	})
}

// 注意：在基于Casbin的设计中，菜单通过permission_key字段关联权限，不需要关联表
