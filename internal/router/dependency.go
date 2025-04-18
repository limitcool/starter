package router

import (
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/storage"
	"github.com/limitcool/starter/internal/repository"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/internal/storage/database"
	"gorm.io/gorm"
)

// Repositories 仓库层依赖集合
type Repositories struct {
	MenuRepo         *repository.MenuRepo
	RoleRepo         *repository.RoleRepo
	SysUserRepo      *repository.SysUserRepo
	PermissionRepo   *repository.PermissionRepo
	OperationLogRepo *repository.OperationLogRepo
	UserRepo         *repository.UserRepo
	FileRepo         *repository.FileRepo
	APIRepo          *repository.APIRepo
}

// Services 服务层依赖集合
type Services struct {
	CasbinService       casbin.Service
	RoleService         *services.RoleService
	MenuService         *services.MenuService
	SysUserService      *services.SysUserService
	PermissionService   *services.PermissionService
	OperationLogService *services.OperationLogService
	SystemService       *services.SystemService
	UserService         *services.UserService
	APIService          *services.APIService
	MenuAPIService      *services.MenuAPIService
}

// Controllers 控制器层依赖集合
type Controllers struct {
	UserController         *controller.UserController
	SysUserController      *controller.SysUserController
	RoleController         *controller.RoleController
	MenuController         *controller.MenuController
	PermissionController   *controller.PermissionController
	OperationLogController *controller.OperationLogController
	SystemController       *controller.SystemController
	FileController         *controller.FileController
	APIController          *controller.APIController
}

// initRepositories 初始化仓库层
func initRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		MenuRepo:         repository.NewMenuRepo(db),
		RoleRepo:         repository.NewRoleRepo(db),
		SysUserRepo:      repository.NewSysUserRepo(db),
		PermissionRepo:   repository.NewPermissionRepo(db),
		OperationLogRepo: repository.NewOperationLogRepo(db),
		UserRepo:         repository.NewUserRepo(db),
		FileRepo:         repository.NewFileRepo(db),
		APIRepo:          repository.NewAPIRepo(db),
	}
}

// initServices 初始化服务层
func initServices(repos *Repositories, casbinService casbin.Service, db database.DB) *Services {
	// 创建服务
	svcs := &Services{
		CasbinService:       casbinService,
		RoleService:         services.NewRoleService(repos.RoleRepo, casbinService),
		MenuService:         services.NewMenuService(repos.MenuRepo, casbinService),
		PermissionService:   services.NewPermissionService(repos.PermissionRepo, repos.RoleRepo, repos.MenuRepo, casbinService),
		OperationLogService: services.NewOperationLogService(repos.OperationLogRepo),
		SystemService:       services.NewSystemService(db),
		UserService:         services.NewUserService(repos.UserRepo),
		APIService:          services.NewAPIService(repos.APIRepo),
	}

	// 创建 SysUserService 并设置 RoleService
	svcs.SysUserService = services.NewSysUserService(repos.SysUserRepo, repos.UserRepo, svcs.RoleService, casbinService)

	// 创建 MenuAPIService
	svcs.MenuAPIService = services.NewMenuAPIService(
		repos.MenuRepo,
		repos.APIRepo,
		repos.RoleRepo,
		repos.PermissionRepo,
		casbinService,
	)

	return svcs
}

// initControllers 初始化控制器层
func initControllers(svcs *Services, repos *Repositories, stg *storage.Storage) *Controllers {
	// 创建控制器
	controllers := &Controllers{
		UserController:         controller.NewUserController(svcs.SysUserService, svcs.UserService),
		SysUserController:      controller.NewSysUserController(svcs.SysUserService),
		RoleController:         controller.NewRoleController(svcs.RoleService, svcs.MenuService),
		MenuController:         controller.NewMenuController(svcs.MenuService),
		PermissionController:   controller.NewPermissionController(svcs.PermissionService),
		OperationLogController: controller.NewOperationLogController(svcs.OperationLogService),
		SystemController:       controller.NewSystemController(svcs.SystemService),
		APIController:          controller.NewAPIController(svcs.APIService, svcs.MenuAPIService),
	}

	// 如果存储组件可用，创建文件控制器
	if stg != nil {
		controllers.FileController = controller.NewFileController(stg, repos.FileRepo)
	}

	return controllers
}
