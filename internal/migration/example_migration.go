package migration

import (
	"gorm.io/gorm"
)

// 这是一个示例迁移文件，展示如何添加新的迁移
// 在实际项目中，你可以创建类似的文件来添加新的迁移

func init() {
	// 注册一个示例迁移
	// 注意：迁移名称应该是唯一的，并且按照顺序排列
	// 这里使用 "999_" 前缀，确保它在所有其他迁移之后执行
	RegisterMigration("999_example_migration", exampleMigrationUp, exampleMigrationDown)
}

// 示例迁移的向上函数
func exampleMigrationUp(tx *gorm.DB) error {
	// 这里可以执行任何数据库操作
	// 例如，创建一个新表
	type ExampleTable struct {
		ID          uint   `gorm:"primarykey"`
		Name        string `gorm:"size:100;not null"`
		Description string `gorm:"size:255"`
	}

	// 注意：在实际项目中，你应该在 model 包中定义模型
	// 这里只是为了示例
	if err := tx.AutoMigrate(&ExampleTable{}); err != nil {
		return err
	}

	// 你也可以执行原始 SQL
	// if err := tx.Exec("CREATE INDEX idx_example_name ON example_tables(name)").Error; err != nil {
	//     return err
	// }

	// 或者插入初始数据
	// if err := tx.Create(&ExampleTable{Name: "Example", Description: "This is an example"}).Error; err != nil {
	//     return err
	// }

	return nil
}

// 示例迁移的向下函数
func exampleMigrationDown(tx *gorm.DB) error {
	// 这里应该撤销向上函数所做的更改
	// 例如，删除创建的表
	return tx.Migrator().DropTable("example_tables")
}

// 如何添加新的业务迁移
// 1. 创建一个新的文件，命名为 XXX_migration.go
// 2. 在 init 函数中注册迁移
// 3. 实现向上和向下函数

/*
例如，如果你想添加一个产品表迁移，可以创建一个 product_migration.go 文件：

```go
package migration

import (
	"github.com/limitcool/starter/internal/model"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("011_create_products_table", createProductsTable, dropProductsTable)
}

func createProductsTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.Product{})
}

func dropProductsTable(tx *gorm.DB) error {
	return tx.Migrator().DropTable("products")
}
```

然后在 model 包中定义 Product 模型：

```go
package model

type Product struct {
	BaseModel
	Name        string  `gorm:"size:100;not null;comment:产品名称"`
	Description string  `gorm:"size:255;comment:产品描述"`
	Price       float64 `gorm:"type:decimal(10,2);not null;comment:产品价格"`
	Stock       int     `gorm:"not null;default:0;comment:库存"`
	Status      int     `gorm:"not null;default:1;comment:状态 1:上架 0:下架"`
}

func (Product) TableName() string {
	return "products"
}
```
*/
