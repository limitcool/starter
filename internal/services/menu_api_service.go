package services

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/repository"
)

// MenuAPIService 菜单API关联服务
type MenuAPIService struct {
	menuRepo       *repository.MenuRepo
	apiRepo        *repository.APIRepo
	roleRepo       *repository.RoleRepo
	permissionRepo *repository.PermissionRepo
	casbinService  casbin.Service
}

// NewMenuAPIService 创建菜单API关联服务
func NewMenuAPIService(
	menuRepo *repository.MenuRepo,
	apiRepo *repository.APIRepo,
	roleRepo *repository.RoleRepo,
	permissionRepo *repository.PermissionRepo,
	casbinService casbin.Service,
) *MenuAPIService {
	return &MenuAPIService{
		menuRepo:       menuRepo,
		apiRepo:        apiRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		casbinService:  casbinService,
	}
}

// AssignAPIsToMenu 为菜单分配API
// 将API关联到菜单，并自动创建权限记录，同时将权限同步到Casbin
func (s *MenuAPIService) AssignAPIsToMenu(menuID uint, apiIDs []uint) error {
	// 获取菜单
	menu, err := s.menuRepo.GetByID(menuID)
	if err != nil {
		return fmt.Errorf("获取菜单失败: %w", err)
	}

	// 开始数据库事务
	tx := s.menuRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务中清除现有关联
	if err := tx.Where("menu_id = ?", menuID).Delete(&model.MenuAPI{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("清除菜单API关联失败: %w", err)
	}

	// 删除菜单相关的权限记录
	if err := tx.Where("menu_id = ? AND api_id IS NOT NULL", menuID).Delete(&model.Permission{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除菜单API权限失败: %w", err)
	}

	// 准备批量插入的数据
	var menuAPIs []model.MenuAPI
	var permissions []model.Permission
	var casbinPolicies [][]string

	// 获取拥有该菜单的所有角色
	roles, err := s.roleRepo.GetRolesByMenuID(menuID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("获取菜单角色失败: %w", err)
	}

	// 准备数据
	for _, apiID := range apiIDs {
		// 获取API
		api, err := s.apiRepo.GetByID(apiID)
		if err != nil {
			log.Warn("获取API失败", "api_id", apiID, "error", err)
			continue
		}

		// 添加菜单API关联
		menuAPIs = append(menuAPIs, model.MenuAPI{
			MenuID: menuID,
			APIID:  apiID,
		})

		// 创建权限记录
		permCode := fmt.Sprintf("%s:%s", api.Path, api.Method)
		permissions = append(permissions, model.Permission{
			Name:        fmt.Sprintf("%s-%s", menu.Name, api.Description),
			Code:        permCode,
			Type:        2, // API类型
			Description: fmt.Sprintf("菜单[%s]关联的API[%s]", menu.Name, api.Path),
			Enabled:     true,
			MenuID:      menuID,
			APIID:       apiID,
		})

		// 为每个角色准备Casbin策略
		for _, role := range roles {
			casbinPolicies = append(casbinPolicies, []string{role.Code, api.Path, api.Method})
		}
	}

	// 批量插入菜单API关联
	if len(menuAPIs) > 0 {
		if err := tx.Create(&menuAPIs).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("批量创建菜单API关联失败: %w", err)
		}
	}

	// 批量插入权限记录
	if len(permissions) > 0 {
		if err := tx.Create(&permissions).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("批量创建权限记录失败: %w", err)
		}
	}

	// 提交数据库事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	// 批量添加Casbin策略
	if len(casbinPolicies) > 0 {
		success, err := s.casbinService.AddPolicies(casbinPolicies)
		if err != nil || !success {
			log.Error("批量添加Casbin权限策略失败", "error", err)
			// 这里不回滚数据库事务，因为数据库事务已经提交
			// 可以在后续的同步操作中修复Casbin状态
		} else {
			log.Info("批量添加Casbin权限策略成功", "count", len(casbinPolicies))
		}
	}

	return nil
}

// AssignMenuToRole 为角色分配菜单时，同步分配API权限
// 将菜单关联到角色，并自动同步菜单关联的API权限到Casbin
func (s *MenuAPIService) AssignMenuToRole(roleID uint, menuIDs []uint) error {
	// 获取角色
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return fmt.Errorf("获取角色失败: %w", err)
	}

	// 开始数据库事务
	tx := s.menuRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务中删除原有的角色菜单关联
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除角色菜单关联失败: %w", err)
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
		if err := tx.Create(&roleMenus).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建角色菜单关联失败: %w", err)
		}
	}

	// 提交数据库事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	// 准备批量添加的Casbin策略
	var casbinPolicies [][]string

	// 为每个菜单同步API权限
	for _, menuID := range menuIDs {
		// 获取菜单关联的所有API
		apis, err := s.apiRepo.GetByMenuID(menuID)
		if err != nil {
			log.Warn("获取菜单API失败", "menu_id", menuID, "error", err)
			continue
		}

		// 为每个API准备Casbin策略
		for _, api := range apis {
			casbinPolicies = append(casbinPolicies, []string{role.Code, api.Path, api.Method})
		}
	}

	// 批量添加Casbin策略
	if len(casbinPolicies) > 0 {
		// 先删除该角色的所有权限
		_, err := s.casbinService.DeleteRole(role.Code)
		if err != nil {
			log.Error("删除角色权限失败", "role", role.Code, "error", err)
			// 这里不回滚数据库事务，因为数据库事务已经提交
		}

		// 批量添加策略
		success, err := s.casbinService.AddPolicies(casbinPolicies)
		if err != nil || !success {
			log.Error("批量添加Casbin权限策略失败", "role", role.Code, "error", err)
			// 这里不回滚数据库事务，因为数据库事务已经提交
		} else {
			log.Info("批量添加角色API权限成功", "role", role.Code, "count", len(casbinPolicies))
		}
	}

	return nil
}

// SyncMenuAPIPermissions 同步所有菜单API权限到Casbin
// 遍历所有菜单和角色，将菜单关联的API权限同步到Casbin
func (s *MenuAPIService) SyncMenuAPIPermissions() error {
	// 获取所有菜单
	menus, err := s.menuRepo.GetAll()
	if err != nil {
		return fmt.Errorf("获取所有菜单失败: %w", err)
	}

	// 开始数据库事务
	tx := s.menuRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 清除所有权限记录中的API关联
	if err := tx.Model(&model.Permission{}).Where("api_id IS NOT NULL").Update("api_id", nil).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("清除权限API关联失败: %w", err)
	}

	// 准备批量插入的数据
	var permissions []model.Permission
	var casbinPolicies [][]string

	// 遍历所有菜单
	for _, menu := range menus {
		// 获取菜单关联的所有API
		apis, err := s.apiRepo.GetByMenuID(menu.ID)
		if err != nil {
			log.Warn("获取菜单API失败", "menu_id", menu.ID, "error", err)
			continue
		}

		// 获取拥有该菜单的所有角色
		roles, err := s.roleRepo.GetRolesByMenuID(menu.ID)
		if err != nil {
			log.Warn("获取菜单角色失败", "menu_id", menu.ID, "error", err)
			continue
		}

		// 准备数据
		for _, api := range apis {
			// 创建权限记录
			permCode := fmt.Sprintf("%s:%s", api.Path, api.Method)
			permissions = append(permissions, model.Permission{
				Name:        fmt.Sprintf("%s-%s", menu.Name, api.Description),
				Code:        permCode,
				Type:        2, // API类型
				Description: fmt.Sprintf("菜单[%s]关联的API[%s]", menu.Name, api.Path),
				Enabled:     true,
				MenuID:      menu.ID,
				APIID:       api.ID,
			})

			// 为每个角色准备Casbin策略
			for _, role := range roles {
				casbinPolicies = append(casbinPolicies, []string{role.Code, api.Path, api.Method})
			}
		}
	}

	// 批量插入权限记录
	if len(permissions) > 0 {
		// 先删除所有API相关的权限记录
		if err := tx.Where("api_id IS NOT NULL").Delete(&model.Permission{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("删除API权限记录失败: %w", err)
		}

		// 创建新的权限记录
		if err := tx.Create(&permissions).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("批量创建权限记录失败: %w", err)
		}
	}

	// 提交数据库事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	// 清除所有Casbin策略并重新添加
	if len(casbinPolicies) > 0 {
		// 先清除所有策略
		s.casbinService.GetEnforcer().ClearPolicy()
		// 重新加载策略
		if err := s.casbinService.GetEnforcer().LoadPolicy(); err != nil {
			log.Error("加载Casbin策略失败", "error", err)
			return fmt.Errorf("加载Casbin策略失败: %w", err)
		}

		// 批量添加策略
		success, err := s.casbinService.AddPolicies(casbinPolicies)
		if err != nil || !success {
			log.Error("批量添加Casbin权限策略失败", "error", err)
			return fmt.Errorf("批量添加Casbin权限策略失败: %w", err)
		}

		log.Info("批量添加Casbin权限策略成功", "count", len(casbinPolicies))
	}

	log.Info("同步菜单API权限完成")
	return nil
}

// GetMenuAPIs 获取菜单关联的所有API
// 返回菜单关联的所有API列表
func (s *MenuAPIService) GetMenuAPIs(menuID uint) ([]*model.API, error) {
	return s.apiRepo.GetByMenuID(menuID)
}

// GetAPIRoles 获取API关联的所有角色
// 返回API关联的所有菜单的角色列表
func (s *MenuAPIService) GetAPIRoles(apiID uint) ([]*model.Role, error) {
	// 检查API是否存在
	if _, err := s.apiRepo.GetByID(apiID); err != nil {
		return nil, fmt.Errorf("获取API失败: %w", err)
	}

	// 获取关联该API的所有菜单
	menuIDs, err := s.menuRepo.GetMenuIDsByAPIID(apiID)
	if err != nil {
		return nil, fmt.Errorf("获取API关联的菜单失败: %w", err)
	}

	// 获取拥有这些菜单的所有角色
	var allRoles []*model.Role
	for _, menuID := range menuIDs {
		roles, err := s.roleRepo.GetRolesByMenuID(menuID)
		if err != nil {
			log.Warn("获取菜单角色失败", "menu_id", menuID, "error", err)
			continue
		}

		// 合并角色列表，避免重复
		roleMap := make(map[uint]*model.Role)
		for _, role := range roles {
			roleMap[role.ID] = role
		}

		for _, role := range roleMap {
			allRoles = append(allRoles, role)
		}
	}

	return allRoles, nil
}
