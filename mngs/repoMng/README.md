repoMng - 多库 + 读写分离 + 通用 Repo[T]

提供一个轻量的存储管理器：
- 命名多库：default/audit/report ...
- 读写分离：在 gorm DB 初始化时挂 dbresolver 即可
- 通用仓储：RepoOf[T](db) + Options（条件/排序/分页/强一致读）
- 事务：InTx(name, ctx, fn) 或手动 db.Transaction

快速开始

```go
import (
    "context"
    "github.com/wiidz/goutil/mngs/repoMng"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// 初始化多个库
mgr := repoMng.NewManager()
mainDB, _ := gorm.Open(postgres.Open(mainDSN))
auditDB, _ := gorm.Open(postgres.Open(auditDSN))

// 如需读写分离：请在各自 DB 上注册 dbresolver（Sources/Replicas）
// _ = mainDB.Use(dbresolver.Register(...))

mgr.SetupDefault(mainDB)
mgr.Register("audit", auditDB)

// 默认库 - 通用仓储
type User struct { ID uint64; LoginID, Nickname string }
userRepo := repoMng.RepoOf[User](mgr.Default().DB())

u, _ := userRepo.First(ctx, repoMng.WithEq("login_id", "alice"))
items, total, _ := userRepo.List(ctx,
    repoMng.WithLike("nickname", "ali"),
    repoMng.WithOrder("id DESC"),
    repoMng.WithPage(1, 20),
)

// 强一致读（非事务）
u, _ = userRepo.First(ctx,
    repoMng.WithEq("login_id", "alice"),
    repoMng.WithWriteRoute(),
)

// 指定库 + 事务
_ = mgr.InTx("audit", ctx, func(ctx context.Context, s *repoMng.Set) error {
    logRepo := repoMng.RepoOf[AuditLog](s.DB())
    return logRepo.Create(ctx, &AuditLog{/* ... */})
})
```

API 摘要

- Manager
  - NewManager()
  - SetupDefault(db *gorm.DB) / Register(name string, db *gorm.DB)
  - Default() *Set / For(name string) *Set
  - InTx(name string, ctx, fn(ctx,*Set) error) error
- Set
  - DB() *gorm.DB、BindTx(tx *gorm.DB) *Set
- Repo
  - RepoOf[T](db *gorm.DB) *Repo[T]
  - GetByID / First / List / Create / Update / Delete
- Options
  - WithEq / WithIn / WithLike / WithOrder / WithPage
  - WithSelect / WithPreload / WithScopes
  - WithWriteRoute()（强一致读）

读写分离说明

- 写：默认走主库；读：默认走从库（dbresolver 控制）
- 事务内所有读写均走主库
- 非事务要求强一致读：在查询时加 WithWriteRoute()

多库说明

- 通过 Register(name, db) 注册多个库
- 使用时 mgr.For("name").DB() 获取对应库 *gorm.DB
- 不同库的读写分离配置分别挂载


