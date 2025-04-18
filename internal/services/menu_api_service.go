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

	// 清除现有关联
	if err := s.menuRepo.ClearMenuAPIs(menuID); err != nil {
		return fmt.Errorf("清除菜单API关联失败: %w", err)
	}

	// 创建新关联
	for _, apiID := range apiIDs {
		// 获取API
		api, err := s.apiRepo.GetByID(apiID)
		if err != nil {
			log.Warn("获取API失败", "api_id", apiID, "error", err)
			continue
		}

		// 创建菜单API关联
		if err := s.menuRepo.AddMenuAPI(menuID, apiID); err != nil {
			log.Warn("添加菜单API关联失败", "menu_id", menuID, "api_id", apiID, "error", err)
			continue
		}

		// 创建权限记录
		permCode := fmt.Sprintf("%s:%s", api.Path, api.Method)
		perm := &model.Permission{
			Name:        fmt.Sprintf("%s-%s", menu.Name, api.Description),
			Code:        permCode,
			Type:        2, // API类型
			Description: fmt.Sprintf("菜单[%s]关联的API[%s]", menu.Name, api.Path),
			Enabled:     true,
			MenuID:      menuID,
			APIID:       apiID,
		}

		// 保存权限
		if err := s.permissionRepo.Create(perm); err != nil {
			log.Warn("创建权限记录失败", "menu_id", menuID, "api_id", apiID, "error", err)
			continue
		}

		// 获取拥有该菜单的所有角色
		roles, err := s.roleRepo.GetRolesByMenuID(menuID)
		if err != nil {
			log.Warn("获取菜单角色失败", "menu_id", menuID, "error", err)
			continue
		}

		// 为每个角色添加API权限
		for _, role := range roles {
			// 在Casbin中添加权限策略
			_, err := s.casbinService.AddPermissionForRole(role.Code, api.Path, api.Method)
			if err != nil {
				log.Warn("添加Casbin权限策略失败", "role", role.Code, "api", api.Path, "method", api.Method, "error", err)
				continue
			}

			log.Info("为角色添加API权限成功", "role", role.Code, "api", api.Path, "method", api.Method)
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

	// 先调用原有的菜单分配方法
	if err := s.menuRepo.AssignMenuToRole(roleID, menuIDs); err != nil {
		return err
	}

	// 为每个菜单同步API权限
	for _, menuID := range menuIDs {
		// 获取菜单关联的所有API
		apis, err := s.apiRepo.GetByMenuID(menuID)
		if err != nil {
			log.Warn("获取菜单API失败", "menu_id", menuID, "error", err)
			continue
		}

		// 为每个API添加Casbin权限策略
		for _, api := range apis {
			_, err := s.casbinService.AddPermissionForRole(role.Code, api.Path, api.Method)
			if err != nil {
				log.Warn("添加Casbin权限策略失败", "role", role.Code, "api", api.Path, "method", api.Method, "error", err)
				continue
			}

			log.Info("为角色添加API权限成功", "role", role.Code, "api", api.Path, "method", api.Method)
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

		// 为每个角色添加API权限
		for _, role := range roles {
			for _, api := range apis {
				// 在Casbin中添加权限策略
				_, err := s.casbinService.AddPermissionForRole(role.Code, api.Path, api.Method)
				if err != nil {
					log.Warn("添加Casbin权限策略失败", "role", role.Code, "api", api.Path, "method", api.Method, "error", err)
					continue
				}
			}
		}
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
