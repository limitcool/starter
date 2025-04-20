package model

// MenuTree 菜单树结构
type MenuTree struct {
	*Menu
	Children []*MenuTree `json:"children"`
}
