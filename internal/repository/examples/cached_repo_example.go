package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/cache"
	"github.com/limitcool/starter/internal/repository"
	"gorm.io/gorm"
)

// 示例：如何使用缓存仓库
func CachedRepoExample(db *gorm.DB) {
	// 创建缓存工厂
	cacheFactory := cache.GetFactory()

	// 创建内存缓存
	userCache, err := cacheFactory.Create("user_cache", cache.Memory,
		cache.WithExpiration(5*time.Minute),
		cache.WithMaxEntries(1000),
	)
	if err != nil {
		panic(err)
	}

	// 创建基础用户仓库
	baseUserRepo := repository.NewUserRepo(db)

	// 创建缓存用户仓库
	cachedUserRepo := repository.NewCachedRepo[model.User](baseUserRepo.GenericRepo, userCache, "user", 5*time.Minute)

	// 使用缓存仓库
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Username: "test_user",
		Password: "password",
		Email:    "test@example.com",
	}

	// 创建用户（会自动缓存）
	if err := cachedUserRepo.Create(ctx, user); err != nil {
		panic(err)
	}

	fmt.Println("用户创建成功:", user.ID)

	// 获取用户（从缓存获取）
	cachedUser, err := cachedUserRepo.GetByID(ctx, user.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println("从缓存获取用户:", cachedUser.Username)

	// 更新用户（会自动更新缓存）
	cachedUser.Email = "updated@example.com"
	if err := cachedUserRepo.Update(ctx, cachedUser); err != nil {
		panic(err)
	}

	fmt.Println("用户更新成功")

	// 再次获取用户（从缓存获取更新后的数据）
	updatedUser, err := cachedUserRepo.GetByID(ctx, user.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println("从缓存获取更新后的用户:", updatedUser.Email)

	// 使用高级查询
	users, err := cachedUserRepo.FindWithLike(ctx, "username", "test")
	if err != nil {
		panic(err)
	}

	fmt.Printf("模糊查询到 %d 个用户\n", len(users))

	// 使用聚合查询
	count, err := cachedUserRepo.AggregateField(ctx, repository.Count, "id", "username LIKE ?", "%test%")
	if err != nil {
		panic(err)
	}

	fmt.Printf("聚合查询结果: %f\n", count)

	// 使用事务
	err = cachedUserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// 在事务中使用仓库
		txRepo := cachedUserRepo.WithTx(tx)

		// 更新用户
		updatedUser.Nickname = "Transaction Test"
		return txRepo.Update(ctx, updatedUser)
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("事务执行成功")

	// 删除用户（会自动删除缓存）
	if err := cachedUserRepo.Delete(ctx, user.ID); err != nil {
		panic(err)
	}

	fmt.Println("用户删除成功")
}
