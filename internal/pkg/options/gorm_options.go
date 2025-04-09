package options

import (
	"gorm.io/gorm"
)

// Option GORM查询选项
type Option func(*gorm.DB) *gorm.DB

// Apply 应用多个Option到GORM查询
func Apply(db *gorm.DB, opts ...Option) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

// WithPage 分页选项
func WithPage(page, pageSize int) Option {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 20
		} else if pageSize > 100 {
			// 限制最大每页条数为100，防止请求过大数据
			pageSize = 100
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// WithOrder 排序选项
func WithOrder(field string, direction string) Option {
	return func(db *gorm.DB) *gorm.DB {
		// 默认按ID降序
		if field == "" {
			field = "id"
		}

		// 只接受 asc 或 desc
		if direction != "asc" && direction != "desc" {
			direction = "desc" // 默认降序
		}

		return db.Order(field + " " + direction)
	}
}

// WithPreload 预加载关联选项
func WithPreload(relation string, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(relation, args...)
	}
}

// WithJoin 连接查询选项
func WithJoin(query string, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(query, args...)
	}
}

// WithSelect 指定查询字段选项
func WithSelect(query interface{}, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(query, args...)
	}
}

// WithGroup 分组选项
func WithGroup(query string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Group(query)
	}
}

// WithHaving HAVING条件选项
func WithHaving(query interface{}, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Having(query, args...)
	}
}

// WithWhere WHERE条件选项
func WithWhere(query interface{}, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

// WithOrWhere OR WHERE条件选项
func WithOrWhere(query interface{}, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Or(query, args...)
	}
}

// WithLike LIKE条件选项
func WithLike(field string, value string) Option {
	return func(db *gorm.DB) *gorm.DB {
		if value == "" {
			return db
		}
		return db.Where(field+" LIKE ?", "%"+value+"%")
	}
}

// WithExactMatch 精确匹配条件选项
func WithExactMatch(field string, value interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		// 如果值为空，则不添加此条件
		if value == nil || value == "" {
			return db
		}
		return db.Where(field+" = ?", value)
	}
}

// WithTimeRange 时间范围条件选项
func WithTimeRange(field string, start, end interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		if start != nil {
			db = db.Where(field+" >= ?", start)
		}
		if end != nil {
			db = db.Where(field+" <= ?", end)
		}
		return db
	}
}

// WithKeyword 关键字搜索选项
// fields: 要搜索的字段列表
func WithKeyword(keyword string, fields ...string) Option {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" || len(fields) == 0 {
			return db
		}

		query := db
		for i, field := range fields {
			if i == 0 {
				query = query.Where(field+" LIKE ?", "%"+keyword+"%")
			} else {
				query = query.Or(field+" LIKE ?", "%"+keyword+"%")
			}
		}

		return query
	}
}

// WithBaseQuery 应用基础查询条件
func WithBaseQuery(tableName string, status *int, keyword string, keywordFields []string, createBy string, startTime, endTime interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		// 应用状态过滤
		if status != nil {
			db = db.Where(tableName+".status = ?", *status)
		}

		// 应用关键字搜索
		if keyword != "" && len(keywordFields) > 0 {
			db = WithKeyword(keyword, keywordFields...)(db)
		}

		// 应用创建人过滤
		if createBy != "" {
			db = db.Where(tableName+".create_by = ?", createBy)
		}

		// 应用时间范围过滤
		if startTime != nil || endTime != nil {
			db = WithTimeRange(tableName+".created_at", startTime, endTime)(db)
		}

		return db
	}
}
