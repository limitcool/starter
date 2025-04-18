package model

// MenuTree 菜单树结构
type MenuTree struct {
	*Menu
	Children []*MenuTree `json:"children"`
}

// BuildMenuTree 构建菜单树
func BuildMenuTree(menus []*Menu) []*MenuTree {
	// 创建一个映射，用于快速查找菜单
	menuMap := make(map[uint]*MenuTree)
	
	// 将所有菜单转换为树节点
	for _, menu := range menus {
		menuMap[menu.ID] = &MenuTree{
			Menu:     menu,
			Children: []*MenuTree{},
		}
	}
	
	// 构建树结构
	var rootMenus []*MenuTree
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
