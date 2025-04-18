package casbin

// Permission 表示一个权限项
type Permission struct {
	Role   string
	Object string
	Action string
}
