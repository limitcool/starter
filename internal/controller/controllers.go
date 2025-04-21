package controller

// Controllers 全局控制器实例
var Controllers = &ControllerInstances{}

// ControllerInstances 控制器实例
type ControllerInstances struct {
	UserController         *UserController
	AdminUserController    *AdminUserController
	RoleController         *RoleController
	MenuController         *MenuController
	PermissionController   *PermissionController
	OperationLogController *OperationLogController
	FileController         *FileController
	APIController          *APIController
	AdminController        *AdminController
}
